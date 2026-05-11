package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4"
)

// Declaración de variables privadas

// Declaración de variables públicas

var lock = sync.RWMutex{}

type ConnData struct {
	Conn   *pgx.Conn
	ConnID string
}

var connChannels map[string]chan ConnData = make(map[string]chan ConnData)

func init() {

}

/*
    type Connections struct {
        availableConnections map[string]string
        currentConnections   map[string]*pgx.Conn
    }

var connections Connections = Connections{availableConnections: make(map[string]string), currentConnections: make(map[string]*pgx.Conn)}

// var connections map[string]map[string]*pgx.Conn

    func readCurrentConnections(connData *ConnData) (*pgx.Conn, bool) {
        //lock.RLock()
        //defer lock.RUnlock()
        var conn *pgx.Conn
        var found bool
        conn, found = connections.currentConnections[connData]
        return conn, found
    }

    func lenCurrentConnections() int {
        //lock.RLock()
        //defer lock.RUnlock()
        return len(connections.currentConnections)
    }

    func writeCurrentConnections(connData *ConnData, value *pgx.Conn) {
        lock.Lock()
        defer lock.Unlock()
        connections.currentConnections[connData] = value
    }

    func deletedAvailableConnections(connData *ConnData) {
        //lock.Lock()
        //defer lock.Unlock()
        delete(connections.availableConnections, connData)
    }

    func deleteCurrentConnections(connData *ConnData) {
        lock.Lock()
        defer lock.Unlock()
        delete(connections.currentConnections, connData)
    }

    func writeAvailableConnections(connData *ConnData, value string) {
        //lock.Lock()
        //defer lock.Unlock()
        connections.availableConnections[connData] = value
    }

    func lenAvailableConnections() int {
        //lock.RLock()
        //defer lock.RUnlock()
        return len(connections.availableConnections)
    }
*/

func GetConnection(connData *ConnData, clientConfig *DBClientConfig, serverConfig *DBServerConfig) (*ConnData, error) {

	getClientConnection(clientConfig, serverConfig)

	if connData.ConnID == "" {

		if channel, found := readChannel(clientConfig.DatabaseName); found {

			conn := <-channel

			// Verificar que la conexión sigue viva; reconectar si está muerta (broken pipe)
			if pingErr := conn.Conn.Ping(context.Background()); pingErr != nil {
				conn.Conn.Close(context.Background())
				sslmode := clientConfig.SSLMode
				if sslmode == "" {
					sslmode = "disable"
				}
				psqlconn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
					clientConfig.UserName, clientConfig.Password,
					clientConfig.Hostname, clientConfig.Port,
					clientConfig.DatabaseName, sslmode)
				newConn, err := pgx.Connect(context.Background(), psqlconn)
				if err != nil {
					return connData, err
				}
				conn.Conn = newConn
			}

			*connData = conn
		}
	}

	return connData, nil
}

/*
En esta función es cuando se adicionan IDs de conexiones a las conexiones disponibles
*/
func ReleaseConnection(connData *ConnData) error {

	if connData.ConnID != "" {

		if channel, found := readChannel(connData.ConnID); found {

			channel <- *connData
		}
	}

	return nil
}

func readChannel(dbName string) (chan ConnData, bool) {
	channel, found := connChannels[dbName]
	return channel, found
}

func getClientConnection(clientConfig *DBClientConfig, serverConfig *DBServerConfig) error {
	lock.Lock()
	defer lock.Unlock()

	poolSize := 3 // valor por defecto conservador
	if serverConfig != nil && serverConfig.PoolSize > 0 {
		poolSize = int(serverConfig.PoolSize)
	}

	if _, found := readChannel(clientConfig.DatabaseName); !found {

		println("X, ", clientConfig.DatabaseName)
		channel := make(chan ConnData, poolSize)
		connChannels[clientConfig.DatabaseName] = channel
		var wg sync.WaitGroup
		wg.Add(poolSize)
		println("Inicializando conexiones de Postgres...")
		var success bool = true
		for i := 0; i < poolSize; i++ {

			go func() {
				defer wg.Done()
				
				sslmode := clientConfig.SSLMode
				if sslmode == "" {
					sslmode = "disable"
				}
				
				// AQUI SE APLICA LA CONEXIÓN SEGURA SSL OBLIGATORIA PARA RENDER SIN ROMPER LA LÓGICA (configurable)
				var psqlconn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
					clientConfig.UserName, clientConfig.Password, clientConfig.Hostname, clientConfig.Port, clientConfig.DatabaseName, sslmode)

				var conn, err = pgx.Connect(context.Background(), psqlconn)

				if err != nil {
					fmt.Println("PostgresConnection getClientConnection Error: ", err)
					success = false
					return
				}

				err = conn.Ping(context.Background())

				if err != nil {
					fmt.Println("PostgresConnection getClientConnection Ping Error: ", err)
					success = false
					return
				}

				//Se crea una nueva conexión para usar
				channel <- ConnData{Conn: conn, ConnID: clientConfig.DatabaseName}
			}()

		}
		wg.Wait()
		if success {
			println("Inicialización completa!")
		}
	}
	return nil
}

func performTransaction(action string, connData *ConnData, clientConfig *DBClientConfig, serverConfig *DBServerConfig) (*ConnData, error) {
	// Definición de variables

	var err error = nil

	//Obtenemos la conexión
	if connData.ConnID == "" {
		connData, err = GetConnection(connData, clientConfig, serverConfig)
		if err != nil {
			return connData, err
		}
	}

	var query string = ""
	switch action {
	case "BEGIN":
		query = "BEGIN"
	case "COMMIT":
		query = "COMMIT"
	case "ROLLBACK":
		query = "ROLLBACK"
	}

	_, err = connData.Conn.Exec(context.Background(), query)

	if err != nil {
		fmt.Println("PostgresConnection performTransaction Error: ", err)
		return connData, err
	}
	return connData, nil
}

func StartTransaction(connData *ConnData, clientConfig *DBClientConfig, serverConfig *DBServerConfig) (*ConnData, error) {
	return performTransaction("BEGIN", connData, clientConfig, serverConfig)
}
func CommitTransaction(connData *ConnData, clientConfig *DBClientConfig, serverConfig *DBServerConfig) (*ConnData, error) {
	return performTransaction("COMMIT", connData, clientConfig, serverConfig)
}
func RollbackTransaction(connData *ConnData, clientConfig *DBClientConfig, serverConfig *DBServerConfig) (*ConnData, error) {
	return performTransaction("ROLLBACK", connData, clientConfig, serverConfig)
}

// CloseAllConnections cierra todas las conexiones pgx del pool al apagar el servidor.
func CloseAllConnections() {
	lock.Lock()
	defer lock.Unlock()
	for dbName, channel := range connChannels {
		n := len(channel)
		for i := 0; i < n; i++ {
			conn := <-channel
			conn.Conn.Close(context.Background())
		}
		delete(connChannels, dbName)
		fmt.Printf("[pgx] Pool '%s' cerrado (%d conexiones)\n", dbName, n)
	}
}