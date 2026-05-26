/**
 * follow-up-contact-modal
 * Reusable component managing the follow-up call attempt flow:
 * - "Contestó" / "No contestó" decision dialog
 * - Call attempt logging (specifying date, target, and reason)
 * - Rescheduling options (1, 2, or 3 days)
 * - Dynamic closure form injection using <dinamic-form>
 * - Success alert notification
 */
app.component('follow-up-contact-modal', {
    delimiters: ['${', '}'],
    props: {
        currentUser:   { type: String, required: true },
        currentUserId: { type: String, required: true },
        userTeam:      { type: String, required: true }
    },
    data() {
        return {
            modalContesto: null,
            attemptDateTime: '',
            selectedContactTarget: '',
            modalRazon: null,
            modalAcciones: null,
            showClosureFormModal: false,
            showSuccessModal: false,
            closureFormId: 'da8423ab-1a8c-47db-96b7-d10496df571a',
            closureSubmissionId: null,
            activeClosureFollowUp: null,
            reasonText: '',
            loading: false,
            rescheduleLoading: null,
            activeFollowUp: null
        };
    },
    methods: {
        // ── Public API (callable by parent ref) ──────────────────────────

        start(fu) {
            console.log('[FollowUpContactModal] start called for', fu);
            this.activeFollowUp = fu;
            this.initAttemptDateTime();

            if (fu.attempts >= 3) {
                this.modalAcciones = fu;
            } else {
                this.modalContesto = fu;
            }
        },

        suggestClosure(fu) {
            console.log('[FollowUpContactModal] suggestClosure called for', fu);
            this.activeFollowUp = fu;
            this.handleSuggestClosure();
        },

        // ── Internal Helpers & Event Handlers ────────────────────────────

        initAttemptDateTime() {
            const tzoffset = (new Date()).getTimezoneOffset() * 60000;
            this.attemptDateTime = (new Date(Date.now() - tzoffset)).toISOString().slice(0, 16);
        },

        closeModals() {
            this.modalContesto = null;
            this.modalRazon = null;
            this.modalAcciones = null;
            this.showClosureFormModal = false;
            this.showSuccessModal = false;
            this.closureSubmissionId = null;
            this.activeClosureFollowUp = null;
            this.reasonText = '';
            this.selectedContactTarget = '';
            this.rescheduleLoading = null;
            this.activeFollowUp = null;
        },

        async handleYesAnswer() {
            const fu = this.activeFollowUp;
            if (!fu || this.loading) return;

            this.loading = true;
            try {
                let customCreatedAt = null;
                if (this.attemptDateTime) {
                    const parts = this.attemptDateTime.split(/[-T:]/);
                    if (parts.length >= 5) {
                        const year = parseInt(parts[0], 10);
                        const month = parseInt(parts[1], 10) - 1;
                        const day = parseInt(parts[2], 10);
                        const hours = parseInt(parts[3], 10);
                        const minutes = parseInt(parts[4], 10);
                        customCreatedAt = new Date(year, month, day, hours, minutes).toISOString();
                    }
                }

                const response = await fetch(`/api/v1/follow-ups/${fu.id}/attempts`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        reason: 'Sí contestó',
                        was_answered: true,
                        created_at: customCreatedAt
                    })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || 'Error al registrar contacto');
                }

                // Redirect to the management view of the follow-up
                window.location.href = `/salvia/hacer-seguimiento/${fu.id}`;
            } catch (error) {
                console.error('Error:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: error.message || 'No se pudo registrar el contacto'
                });
                this.loading = false;
            }
        },

        handleNoAnswer() {
            const fu = this.activeFollowUp;
            this.modalContesto = null;
            this.modalRazon = fu;
            this.reasonText = '';
            this.selectedContactTarget = '';
        },

        async handleSaveReason() {
            if (!this.reasonText || this.reasonText.trim() === '' || !this.selectedContactTarget) return;

            const fu = this.activeFollowUp;
            this.loading = true;

            const fullReason = "se intento contactar " + this.selectedContactTarget + " " + this.reasonText.trim();

            try {
                let customCreatedAt = null;
                if (this.attemptDateTime) {
                    const parts = this.attemptDateTime.split(/[-T:]/);
                    if (parts.length >= 5) {
                        const year = parseInt(parts[0], 10);
                        const month = parseInt(parts[1], 10) - 1;
                        const day = parseInt(parts[2], 10);
                        const hours = parseInt(parts[3], 10);
                        const minutes = parseInt(parts[4], 10);
                        customCreatedAt = new Date(year, month, day, hours, minutes).toISOString();
                    }
                }

                const response = await fetch(`/api/v1/follow-ups/${fu.id}/attempts`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ 
                        reason: fullReason,
                        created_at: customCreatedAt
                    })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || 'Error saving attempt');
                }

                const data = await response.json();
                
                // Reactive update on the parent card object (passed by reference)
                fu.attempts = data.attempts;
                if (data.followUp && data.followUp.last_attempt_at) {
                    fu.last_attempt_at = data.followUp.last_attempt_at;
                }

                const newAttempt = {
                    id: `temp-${Date.now()}`,
                    follow_up_id: fu.id,
                    reason: fullReason,
                    was_answered: false,
                    created_at: customCreatedAt || new Date().toISOString(),
                    updated_at: customCreatedAt || new Date().toISOString()
                };

                if (!fu.follow_up_attempts) fu.follow_up_attempts = [];
                fu.follow_up_attempts.unshift(newAttempt);

                this.modalRazon = null;
                this.reasonText = '';
                this.selectedContactTarget = '';

                // Emit attempt completed event in case parent dashboard needs to react
                this.$emit('completed', { type: 'attempt', followUp: fu });

                // Show success overlay
                const successOverlay = document.getElementById('successOverlay');
                if (successOverlay) successOverlay.style.display = 'flex';
                setTimeout(() => {
                    if (successOverlay) successOverlay.style.display = 'none';
                }, 2000);

                if (data.maxAttemptsReached) {
                    this.modalAcciones = fu;
                } else {
                    this.closeModals();
                }
            } catch (error) {
                console.error('Error:', error);
                const failOverlay = document.getElementById('failOverlay');
                const failMsg = document.getElementById('fail_msg');
                if (failMsg) failMsg.textContent = error.message;
                if (failOverlay) failOverlay.style.display = 'flex';
                setTimeout(() => {
                    if (failOverlay) failOverlay.style.display = 'none';
                }, 3000);
            } finally {
                this.loading = false;
            }
        },

        async handleReschedule(daysToAdd) {
            const fu = this.activeFollowUp;
            if (!fu || this.rescheduleLoading !== null) return;

            this.rescheduleLoading = daysToAdd;
            this.loading = true;
            try {
                const targetDate = new Date();
                targetDate.setDate(targetDate.getDate() + daysToAdd);
                const fechaDestino = targetDate.toISOString().split('T')[0];

                const response = await fetch(`/api/v1/seguimientos/${fu.id}/reagendar`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        nueva_fecha: fechaDestino,
                        nueva_hora: fu.scheduledTime || '08:00',
                        motivo: `Reagendamiento automático por límite de intentos alcanzado (+${daysToAdd} días)`
                    })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || 'Error al reagendar');
                }

                // Notify parent that follow-up has been rescheduled so it can remove it from dashboard lists
                this.$emit('completed', { type: 'rescheduled', followUpId: fu.id });

                this.modalAcciones = null;
                this.closeModals();

                // Show success overlay
                const successOverlay = document.getElementById('successOverlay');
                if (successOverlay) successOverlay.style.display = 'flex';
                setTimeout(() => {
                    if (successOverlay) successOverlay.style.display = 'none';
                }, 2000);

            } catch (error) {
                console.error('Error:', error);
                const failOverlay = document.getElementById('failOverlay');
                const failMsg = document.getElementById('fail_msg');
                if (failMsg) failMsg.textContent = error.message;
                if (failOverlay) failOverlay.style.display = 'flex';
                setTimeout(() => {
                    if (failOverlay) failOverlay.style.display = 'none';
                }, 3000);
            } finally {
                this.loading = false;
                this.rescheduleLoading = null;
            }
        },

        handleAddAttempt() {
            const fu = this.activeFollowUp;
            this.modalAcciones = null;
            this.initAttemptDateTime();
            this.modalContesto = fu;
        },

        async handleSuggestClosure() {
            const fu = this.activeFollowUp;
            if (!fu || this.loading) return;

            this.loading = true;
            try {
                const response = await fetch(`/api/v1/follow-ups/${fu.id}/init-closure-form?agent_id=${encodeURIComponent(this.currentUserId)}`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' }
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || 'Error al iniciar formulario de cierre');
                }

                const data = await response.json();
                this.activeClosureFollowUp = fu;
                this.closureSubmissionId = data.submission_id;
                this.showClosureFormModal = true;
                this.modalAcciones = null;
            } catch (error) {
                console.error('Error:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: error.message || 'No se pudo iniciar el formulario de cierre'
                });
            } finally {
                this.loading = false;
            }
        },

        onClosureFormCompleted() {
            console.log('[FollowUpContactModal] onClosureFormCompleted - activeClosureFollowUp:', this.activeClosureFollowUp);
            const closedCaseId = this.activeClosureFollowUp ? (this.activeClosureFollowUp.caseId || this.activeClosureFollowUp.case_id) : null;
            const fuId = this.activeClosureFollowUp ? this.activeClosureFollowUp.id : null;
            
            this.closeModals();
            
            // Notify parent that follow-up closure is finished
            this.$emit('completed', { type: 'closed', followUpId: fuId, caseId: closedCaseId });
            this.showSuccessModal = true;
        },

        shouldShowClosureButton(fu) {
            if (!fu || !fu.follow_up_attempts || fu.follow_up_attempts.length < 9) {
                return false;
            }
            
            const uniqueDays = new Set();
            fu.follow_up_attempts.forEach(attempt => {
                if (attempt.created_at) {
                    const date = attempt.created_at.split('T')[0];
                    uniqueDays.add(date);
                }
            });
            
            return uniqueDays.size >= 3;
        },

        getFutureDateLabel(daysToAdd) {
            const date = new Date();
            date.setDate(date.getDate() + daysToAdd);
            
            if (daysToAdd === 1) return 'Mañana';
            
            const d = date.getDate().toString().padStart(2, '0');
            const m = (date.getMonth() + 1).toString().padStart(2, '0');
            const y = date.getFullYear().toString().substring(2);
            return `${d}/${m}/${y}`;
        }
    },
    template: `
        <div>
            <!-- Modal Contesto -->
            <div v-if="modalContesto" class="ms-modal-overlay" @click.self="closeModals">
                <div class="ms-modal-content">
                    <div style="width: 4rem; height: 4rem; border-radius: 9999px; background-color: #f5f3ff; color: #8b5cf6; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.5rem; font-size: 1.5rem;">
                        <i class="fas fa-phone-alt"></i>
                    </div>
                    <h3 style="font-size: 1.25rem; font-weight: 700; color: #111827; margin-bottom: 0.5rem;">¿Contestó la llamada?</h3>
                    <p style="font-size: 0.875rem; color: #6b7280; margin-bottom: 0.5rem;">Estamos contactando a</p>
                    <p style="font-size: 1rem; font-weight: 700; color: #111827; margin-bottom: 1.5rem;">\${ activeFollowUp ? activeFollowUp.caseName : '' }</p>
                    <div v-if="activeFollowUp && activeFollowUp.attempts > 0"
                        style="background-color: #fee2e2; border: 1px solid #ef4444; border-radius: 0.5rem; padding: 0.75rem; margin-bottom: 1.5rem;">
                        <p style="font-size: 0.8125rem; color: #b91c1c; margin: 0; font-weight: 600;">
                            <i class="fas fa-exclamation-triangle"></i>
                            \${ activeFollowUp.attempts } intento\${ activeFollowUp.attempts > 1 ? 's' : '' } sin respuesta
                        </p>
                    </div>

                    <!-- Selector de Fecha y Hora del Intento -->
                    <div style="margin-bottom: 1.5rem; text-align: left; width: 100%;">
                        <label style="display: block; font-size: 0.8125rem; font-weight: 700; color: #374151; margin-bottom: 0.375rem; min-height: auto;">
                            <i class="far fa-calendar-alt"></i> Fecha y hora del intento
                        </label>
                        <input type="datetime-local" v-model="attemptDateTime"
                            style="width: 100%; padding: 0.75rem 1rem; border-radius: 0.75rem; border: 1px solid #e5e7eb; font-size: 0.875rem; outline: none; transition: all 0.2s; background: white; box-sizing: border-box; font-family: inherit; color: #374151;"
                            @focus="$event.target.style.borderColor='#8b5cf6'; $event.target.style.boxShadow='0 0 0 3px rgba(139,92,246,0.1)'"
                            @blur="$event.target.style.borderColor='#e5e7eb'; $event.target.style.boxShadow='none'">
                    </div>

                    <div style="display: flex; gap: 1rem; width: 100%;">
                        <button @click="handleNoAnswer" :disabled="loading"
                            :style="\`flex: 1; background: white; border: 1px solid #e5e7eb; color: \${loading ? '#9ca3af' : '#4b5563'}; padding: 0.875rem 1.25rem; border-radius: 0.75rem; font-weight: 600; cursor: \${loading ? 'not-allowed' : 'pointer'}; transition: all 0.2s; font-size: 0.875rem;\`">
                            No contestó
                        </button>
                        <button @click="handleYesAnswer" :disabled="loading"
                            :style="\`flex: 1; padding: 0.875rem 1.25rem; border-radius: 0.75rem; font-weight: 600; border: none; cursor: \${loading ? 'not-allowed' : 'pointer'}; transition: all 0.2s; font-size: 0.875rem; background-color: \${loading ? '#d1d5db' : '#8b5cf6'}; color: white;\`"
                            style="display: flex; align-items: center; justify-content: center; gap: 0.5rem;">
                            <template v-if="loading">
                                <i class="fas fa-circle-notch fa-spin"></i> Cargando...
                            </template>
                            <template v-else>
                                Sí contestó <i class="fas fa-check"></i>
                            </template>
                        </button>
                    </div>
                    <button @click="closeModals" :disabled="loading"
                        :style="\`margin-top: 1.5rem; background: none; border: none; font-size: 0.75rem; color: #9ca3af; font-weight: 600; cursor: \${loading ? 'not-allowed' : 'pointer'};\`">
                        CANCELAR
                    </button>
                </div>
            </div>

            <!-- Modal Razón -->
            <div v-if="modalRazon" class="ms-modal-overlay" @click.self="closeModals">
                <div class="ms-modal-content">
                    <div style="width: 4rem; height: 4rem; border-radius: 9999px; background-color: #fff7ed; color: #ea580c; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.5rem; font-size: 1.5rem;">
                        <i class="fas fa-phone-slash"></i>
                    </div>
                    <h3 style="font-size: 1.25rem; font-weight: 700; color: #111827; margin-bottom: 0.5rem;">Intento sin respuesta</h3>
                    <p style="font-size: 0.875rem; color: #6b7280; margin-bottom: 1.5rem;">Registra la razón del intento fallido</p>
                    
                    <!-- Selección de Destinatario del Intento -->
                    <div style="margin-bottom: 1.25rem; text-align: left; width: 100%;">
                        <label style="display: block; font-size: 0.8125rem; font-weight: 700; color: #374151; margin-bottom: 0.5rem; min-height: auto;">
                            ¿A quién se intentó contactar? <span style="color: #dc2626;">*</span>
                        </label>
                        <div style="display: flex; gap: 0.5rem; width: 100%; flex-wrap: wrap;">
                            <button type="button" @click="selectedContactTarget = 'Contacto con la víctima'"
                                :style="\`flex: 1 1 calc(33.333% - 0.5rem); text-align: center; padding: 0.5rem 0.25rem; border-radius: 0.5rem; font-size: 0.75rem; font-weight: 600; border: 1.5px solid \${selectedContactTarget === 'Contacto con la víctima' ? '#8b5cf6' : '#e5e7eb'}; background-color: \${selectedContactTarget === 'Contacto con la víctima' ? '#f5f3ff' : 'white'}; color: \${selectedContactTarget === 'Contacto con la víctima' ? '#7c3aed' : '#4b5563'}; cursor: pointer; transition: all 0.2s;\`"
                                style="box-sizing: border-box; display: flex; align-items: center; justify-content: center; min-height: 2.75rem; line-height: 1.2;">
                                Contacto con la víctima
                            </button>
                            <button type="button" @click="selectedContactTarget = 'Contacto indirecto - Familiar o persona conocida'"
                                :style="\`flex: 1 1 calc(33.333% - 0.5rem); text-align: center; padding: 0.5rem 0.25rem; border-radius: 0.5rem; font-size: 0.75rem; font-weight: 600; border: 1.5px solid \${selectedContactTarget === 'Contacto indirecto - Familiar o persona conocida' ? '#8b5cf6' : '#e5e7eb'}; background-color: \${selectedContactTarget === 'Contacto indirecto - Familiar o persona conocida' ? '#f5f3ff' : 'white'}; color: \${selectedContactTarget === 'Contacto indirecto - Familiar o persona conocida' ? '#7c3aed' : '#4b5563'}; cursor: pointer; transition: all 0.2s;\`"
                                style="box-sizing: border-box; display: flex; align-items: center; justify-content: center; min-height: 2.75rem; line-height: 1.2;">
                                Contacto indirecto - Familiar o persona conocida
                            </button>
                            <button type="button" @click="selectedContactTarget = 'Contacto Con Institución'"
                                :style="\`flex: 1 1 calc(33.333% - 0.5rem); text-align: center; padding: 0.5rem 0.25rem; border-radius: 0.5rem; font-size: 0.75rem; font-weight: 600; border: 1.5px solid \${selectedContactTarget === 'Contacto Con Institución' ? '#8b5cf6' : '#e5e7eb'}; background-color: \${selectedContactTarget === 'Contacto Con Institución' ? '#f5f3ff' : 'white'}; color: \${selectedContactTarget === 'Contacto Con Institución' ? '#7c3aed' : '#4b5563'}; cursor: pointer; transition: all 0.2s;\`"
                                style="box-sizing: border-box; display: flex; align-items: center; justify-content: center; min-height: 2.75rem; line-height: 1.2;">
                                Contacto Con Institución
                            </button>
                        </div>
                    </div>

                    <textarea v-model="reasonText" placeholder="Ej: El usuario tenía el teléfono apagado..." rows="3"
                        style="width: 100%; padding: 0.75rem 1rem; border-radius: 0.75rem; border: 1px solid #e5e7eb; font-size: 0.875rem; outline: none; transition: all 0.2s; background: white; box-sizing: border-box; resize: none; font-family: inherit;"></textarea>
                    <div style="display: flex; gap: 1rem; margin-top: 1.5rem;">
                        <button @click="closeModals"
                            style="flex: 1; background: white; border: 1px solid #e5e7eb; color: #4b5563; padding: 0.875rem 1.25rem; border-radius: 0.75rem; font-weight: 600; cursor: pointer; transition: all 0.2s; font-size: 0.875rem;">
                            Cancelar
                        </button>
                        <button @click="handleSaveReason" :disabled="loading || !reasonText || reasonText.trim() === '' || !selectedContactTarget"
                            :style="\`flex: 1; padding: 0.875rem 1.25rem; border-radius: 0.75rem; font-weight: 600; border: none; cursor: \${loading || !reasonText || reasonText.trim() === '' || !selectedContactTarget ? 'not-allowed' : 'pointer'}; transition: all 0.2s; font-size: 0.875rem; background-color: \${loading || !reasonText || reasonText.trim() === '' || !selectedContactTarget ? '#d1d5db' : '#8b5cf6'}; color: \${loading || !reasonText || reasonText.trim() === '' || !selectedContactTarget ? '#9ca3af' : 'white'};\`"
                            style="display: flex; align-items: center; justify-content: center; gap: 0.5rem;">
                            <template v-if="loading">
                                <i class="fas fa-circle-notch fa-spin"></i> Guardando...
                            </template>
                            <template v-else>
                                Guardar
                            </template>
                        </button>
                    </div>
                </div>
            </div>

            <!-- Modal Acciones -->
            <div v-if="modalAcciones" class="ms-modal-overlay" @click.self="closeModals">
                <div class="ms-modal-content">
                    <div style="width: 4rem; height: 4rem; border-radius: 9999px; background-color: #fee2e2; color: #dc2626; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.5rem; font-size: 1.5rem;">
                        <i class="fas fa-exclamation-circle"></i>
                    </div>
                    <h3 style="font-size: 1.25rem; font-weight: 700; color: #111827; margin-bottom: 0.5rem;">Límite de intentos</h3>
                    <p style="font-size: 0.875rem; color: #6b7280; margin-bottom: 0.5rem;">Se han realizado</p>
                    <p style="font-size: 1.5rem; font-weight: 800; color: #dc2626; margin-bottom: 0.5rem;">\${ modalAcciones.attempts } intentos</p>
                    <p style="font-size: 0.875rem; color: #6b7280; margin-bottom: 2rem;">sin respuesta para contactar a \${ modalAcciones.caseName }</p>
                    <div style="display: flex; flex-direction: column; gap: 0.75rem;">
                        <div style="display: flex; flex-direction: column; gap: 0.5rem; background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 0.75rem; padding: 0.75rem;">
                            <p style="font-size: 0.75rem; font-weight: 700; color: #4b5563; text-align: center; margin: 0 0 0.5rem 0;">POSPONER SEGUIMIENTO PARA:</p>
                            <div style="display: flex; gap: 0.5rem;">
                                <button @click="handleReschedule(1)" :disabled="loading || rescheduleLoading !== null"
                                    :style="\`flex: 1; background: white; border: 1px solid #d1d5db; color: #374151; padding: 0.625rem; border-radius: 0.5rem; font-weight: 600; cursor: \${loading || rescheduleLoading !== null ? 'not-allowed' : 'pointer'}; transition: all 0.2s; font-size: 0.8125rem; opacity: \${loading || (rescheduleLoading !== null && rescheduleLoading !== 1) ? '0.5' : '1'}; display: flex; align-items: center; justify-content: center; gap: 0.375rem;\`">
                                    <i v-if="rescheduleLoading === 1" class="fas fa-circle-notch fa-spin" style="font-size: 0.75rem;"></i>
                                    \${ getFutureDateLabel(1) }
                                </button>
                                <button @click="handleReschedule(2)" :disabled="loading || rescheduleLoading !== null"
                                    :style="\`flex: 1; background: white; border: 1px solid #d1d5db; color: #374151; padding: 0.625rem; border-radius: 0.5rem; font-weight: 600; cursor: \${loading || rescheduleLoading !== null ? 'not-allowed' : 'pointer'}; transition: all 0.2s; font-size: 0.8125rem; opacity: \${loading || (rescheduleLoading !== null && rescheduleLoading !== 2) ? '0.5' : '1'}; display: flex; align-items: center; justify-content: center; gap: 0.375rem;\`">
                                    <i v-if="rescheduleLoading === 2" class="fas fa-circle-notch fa-spin" style="font-size: 0.75rem;"></i>
                                    \${ getFutureDateLabel(2) }
                                </button>
                                <button @click="handleReschedule(3)" :disabled="loading || rescheduleLoading !== null"
                                    :style="\`flex: 1; background: white; border: 1px solid #d1d5db; color: #374151; padding: 0.625rem; border-radius: 0.5rem; font-weight: 600; cursor: \${loading || rescheduleLoading !== null ? 'not-allowed' : 'pointer'}; transition: all 0.2s; font-size: 0.8125rem; opacity: \${loading || (rescheduleLoading !== null && rescheduleLoading !== 3) ? '0.5' : '1'}; display: flex; align-items: center; justify-content: center; gap: 0.375rem;\`">
                                    <i v-if="rescheduleLoading === 3" class="fas fa-circle-notch fa-spin" style="font-size: 0.75rem;"></i>
                                    \${ getFutureDateLabel(3) }
                                </button>
                            </div>
                        </div>
                        <button @click="handleAddAttempt" :disabled="loading || rescheduleLoading !== null"
                            :style="\`background: white; border: 1px solid #e5e7eb; color: #4b5563; padding: 0.875rem 1.25rem; border-radius: 0.75rem; font-weight: 600; cursor: \${loading || rescheduleLoading !== null ? 'not-allowed' : 'pointer'}; transition: all 0.2s; font-size: 0.875rem; width: 100%; display: flex; align-items: center; justify-content: center; gap: 0.5rem; opacity: \${loading || rescheduleLoading !== null ? '0.5' : '1'};\`">
                            <i class="fas fa-plus-circle"></i> Añadir intento de contacto
                        </button>
                        <button v-if="shouldShowClosureButton(modalAcciones)" @click="handleSuggestClosure" :disabled="loading || rescheduleLoading !== null"
                            :style="\`background-color: #dc2626; color: white; padding: 0.875rem 1.25rem; border-radius: 0.75rem; font-weight: 600; border: none; cursor: \${loading || rescheduleLoading !== null ? 'not-allowed' : 'pointer'}; transition: all 0.2s; font-size: 0.875rem; width: 100%; display: flex; align-items: center; justify-content: center; gap: 0.5rem; opacity: \${loading || rescheduleLoading !== null ? '0.5' : '1'};\`">
                            <template v-if="loading">
                                <i class="fas fa-circle-notch fa-spin"></i> Cargando formulario...
                            </template>
                            <template v-else>
                                <i class="fas fa-times-circle"></i> Cerrar caso
                            </template>
                        </button>
                    </div>
                    <button @click="closeModals" :disabled="loading || rescheduleLoading !== null"
                        :style="\`margin-top: 1.5rem; background: none; border: none; font-size: 0.75rem; color: #9ca3af; font-weight: 600; cursor: \${loading || rescheduleLoading !== null ? 'not-allowed' : 'pointer'}; opacity: \${loading || rescheduleLoading !== null ? '0.5' : '1'};\`">
                        VOLVER
                    </button>
                </div>
            </div>

            <!-- Modal Formulario de Cierre Dinámico -->
            <div v-if="showClosureFormModal" class="ms-modal-overlay" @click.self="closeModals">
                <div class="ms-modal-content" style="max-width: 1000px !important; width: 95% !important; max-height: 90vh; overflow: auto !important; display: flex; flex-direction: column; border-radius: 1rem; padding: 1.5rem !important; background: white; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04); text-align: left !important;">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; border-bottom: 1px solid #e5e7eb; padding-bottom: 1rem; flex-shrink: 0;">
                        <div style="display: flex; align-items: center; gap: 0.75rem;">
                            <div style="width: 2.5rem; height: 2.5rem; border-radius: 0.5rem; background-color: #fee2e2; color: #dc2626; display: flex; align-items: center; justify-content: center; font-size: 1.25rem;">
                                <i class="fas fa-exclamation-triangle"></i>
                            </div>
                            <div>
                                <h3 style="font-size: 1.25rem; font-weight: 700; color: #111827; margin: 0;">Formulario de Cierre de Caso</h3>
                                <p style="font-size: 0.875rem; color: #6b7280; margin: 0.25rem 0 0 0;">
                                    Cierre de caso para: <strong style="color: #111827;">\${ activeClosureFollowUp ? activeClosureFollowUp.caseName : '' }</strong>
                                </p>
                            </div>
                        </div>
                        <button @click="closeModals" style="background: none; border: none; font-size: 1.25rem; color: #9ca3af; cursor: pointer; padding: 0.25rem; transition: color 0.2s;" onmouseover="this.style.color='#4b5563'" onmouseout="this.style.color='#9ca3af'">
                             <i class="fas fa-times"></i>
                        </button>
                    </div>
                    
                    <div style="flex: 1; min-height: 0;">
                        <dinamic-form 
                            v-if="closureFormId && closureSubmissionId" 
                            :form-id="closureFormId" 
                            :submission-id="closureSubmissionId" 
                            @form-completed="onClosureFormCompleted">
                        </dinamic-form>
                    </div>
                </div>
            </div>

            <!-- Modal de Éxito al Guardar Formulario -->
            <div v-if="showSuccessModal" class="ms-modal-overlay" @click.self="showSuccessModal = false">
                <div class="ms-modal-content" style="max-width: 28rem !important; padding: 2rem !important; text-align: center !important;">
                    <div style="width: 4rem; height: 4rem; border-radius: 9999px; background-color: #f0fdf4; color: #16a34a; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.5rem; font-size: 2rem;">
                        <i class="fas fa-check-circle"></i>
                    </div>
                    <h3 style="font-size: 1.5rem; font-weight: 700; color: #111827; margin: 0 0 0.5rem 0; min-height: auto; line-height: 1.2;">Formulario Guardado</h3>
                    <p style="font-size: 0.9375rem; color: #4b5563; margin: 0 0 2rem 0; line-height: 1.5;">El formulario de cierre se guardó con éxito.</p>
                    <button @click="showSuccessModal = false" style="width: 100%; padding: 0.875rem 1.25rem; border-radius: 0.75rem; font-weight: 600; border: none; cursor: pointer; transition: all 0.2s; font-size: 0.875rem; background-color: #7c3aed; color: white; outline: none; box-shadow: 0 4px 6px -1px rgba(124, 58, 237, 0.2);" onmouseover="this.style.backgroundColor='#6d28d9'" onmouseout="this.style.backgroundColor='#7c3aed'">
                        OK
                    </button>
                </div>
            </div>
        </div>
    `
});
