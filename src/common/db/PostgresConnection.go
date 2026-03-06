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

	/*func readCurrentConnections(connData *ConnData) (*pgx.Conn, bool) {
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

	getClientConnection(clientConfig)

	//println("Conn status - Curr num conn: "+strconv.Itoa(lenCurrentConnections()), ", Avail num conn: ")

	if connData.ConnID == "" {

		if channel, found := readChannel(clientConfig.DatabaseName); found {

			*connData = <-channel
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
func getClientConnection(clientConfig *DBClientConfig) error {
	lock.Lock()
	defer lock.Unlock()

	/*
		TODO:
			1. Poner variable en PostgresConnection o alguna configuración global con el número máximo de conexiones que aguanta postgres en total
			2. En clientConfig poner el número de conexiones que se le asignarán al susuario en su base de datos
			3. Adicionar en ConnData un campo de tiempo para registrar la última vez que se usó alguna conexión de esa BD
			4. Cada vez que haya una bd nueva, crear su respectivo connChannels[clientConfig.DatabaseName]
				4.1 Si la BD es nueva pero el número de conexiones para ese cliente supera el total de conexiones existentes,
					Entonces se identifica el cliente/bd en connChannels con el tiempo más grande de inactividad, cierran las conexiones y se elimina de connChannels
				4.2 el punto 4.1 se hace en ciclo hasta que hayan conexiones disponibles para darle soporte al nuevo cliente/bd
	*/
	if _, found := readChannel(clientConfig.DatabaseName); !found {

		println("X, ", clientConfig.DatabaseName)
		channel := make(chan ConnData, 80)
		connChannels[clientConfig.DatabaseName] = channel
		var wg sync.WaitGroup
		wg.Add(80)
		println("Inicializando conexiones de Postgres...")
		var success bool = true
		for i := 0; i < 80; i++ {

			go func() {
				var psqlconn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
					clientConfig.UserName, clientConfig.Password, clientConfig.Hostname, clientConfig.Port, clientConfig.DatabaseName)

				var conn, err = pgx.Connect(context.Background(), psqlconn)

				if err != nil {
					fmt.Println("PostgresConnection getClientConnection Error: ", err)
					success = false
				}

				err = conn.Ping(context.Background())

				if err != nil {
					fmt.Println("PostgresConnection getClientConnection Ping Error: ", err)
					success = false
				}

				//Se crea una nueva conexión para usar
				channel <- ConnData{Conn: conn, ConnID: clientConfig.DatabaseName}
				wg.Done()
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
