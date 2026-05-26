# Guía de Integración: Componente Reutilizable `follow-up-contact-modal`

Este componente encapsula todo el flujo modal del contacto de un seguimiento (diálogo ¿Contestó?, registro de intentos con fecha/hora, reprogramación automática y el formulario dinámico de cierre de caso).

---

## 1. Importación del Script

Para poder usar el componente, primero debes incluir su archivo script en el HTML de la página, asegurándote de que se cargue después de Vue y antes de montar la aplicación:

```html
<!-- Importar componente de contacto reutilizable -->
<script src="/static/js/components/follow-up-contact-modal.js"></script>
```

---

## 2. Declaración en el HTML (dentro del contenedor Vue)

Instancia el componente agregando la etiqueta en el HTML de tu página. Debe estar dentro del elemento donde se monta tu aplicación de Vue (por ejemplo, `#app`).

Debes pasarle los datos de la sesión del agente como props (`current-user`, `current-user-id`, `user-team`), y escuchar el evento `@completed`:

```html
<follow-up-contact-modal 
    ref="contactModal"
    :current-user="currentUser" 
    :current-user-id="currentUserId" 
    :user-team="userTeam"
    @completed="onFollowUpContactCompleted">
</follow-up-contact-modal>
```

> [!NOTE]
> **Sobre los Props**:
> Se pasan por props para evitar acoplar el componente a los nombres de variables que usa la página principal. Por ejemplo, en "Detalle del Caso" el ID del usuario se llama `userICode`, por lo que allí lo instanciarías como `:current-user-id="userICode"`.

---

## 3. Disparo del Flujo desde JavaScript

Para abrir los modales de contacto, debes obtener la referencia del componente usando `this.$refs` y llamar al método público `.start(fu)` pasándole el objeto completo del seguimiento actual:

```javascript
// Método en el Vue principal de tu página
iniciarContacto(fu) {
    // 'fu' debe ser un objeto con la estructura básica del seguimiento:
    // {
    //     id: "uuid-del-seguimiento",
    //     caseId: "uuid-del-caso",
    //     caseName: "Nombre de la Víctima",
    //     attempts: 0,
    //     follow_up_attempts: [...]
    // }
    this.$refs.contactModal.start(fu);
}
```

---

## 4. Captura de Eventos de Finalización (`@completed`)

Cuando el componente realiza una acción definitiva sobre el seguimiento (reprogramar por límite de intentos fallidos o realizar el cierre del caso), emite un evento `@completed` con los detalles. Debes manejar este evento en tu Vue principal para actualizar el listado en pantalla:

```javascript
onFollowUpContactCompleted({ type, followUpId, caseId }) {
    console.log('[Dashboard] Evento completado recibido del modal:', type, followUpId, caseId);
    
    if (type === 'rescheduled') {
        // El seguimiento fue pospuesto (ej. remover el card de la lista local de pendientes)
        this.pendingFollowUps = this.pendingFollowUps.filter(f => f.id !== followUpId);
    } else if (type === 'closed') {
        // El caso fue cerrado exitosamente (ej. remover todos los seguimientos asociados a ese caso)
        if (caseId) {
            this.pendingFollowUps = this.pendingFollowUps.filter(f => f.caseId !== caseId);
            this.completedFollowUps = this.completedFollowUps.filter(f => f.caseId !== caseId);
        }
        // Opcional: Recargar los datos después de un breve delay para asegurar la propagación en DB
        setTimeout(async () => {
            await this.loadFollowUps();
        }, 800);
    }
}
```
