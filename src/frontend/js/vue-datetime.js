let _idCounter = 0

const DateTimeInput = {
  name: 'DateTimeInput',
  template: `
    <fieldset class="datetime-input-wrapper" :class="{ 'has-error': errorMessage, 'is-focused': isFocused }" :aria-describedby="descriptionIds">
      <div class="datetime-field">
        <template v-for="(segment, index) in segments" :key="index">
          <span v-if="segment.type === 'separator'" class="datetime-separator" aria-hidden="true">{{ segment.char }}</span>
          <input
            v-else
            :ref="el => segmentRefs[index] = el"
            :id="index === 0 ? inputId : inputId + '-' + index"
            type="text"
            inputmode="numeric"
            class="datetime-segment"
            :class="'segment-' + segment.token"
            :value="segment.display"
            :placeholder="segment.placeholder"
            :maxlength="segment.length"
            :aria-label="index === fieldIndices[0] ? (label ? label + ' - ' + segment.ariaLabel : segment.ariaLabel) : segment.ariaLabel"
            :aria-valuemin="segment.min"
            :aria-valuemax="segment.max"
            :aria-valuenow="segment.numericValue || undefined"
            :aria-invalid="errorMessage ? 'true' : undefined"
            role="spinbutton"
            autocomplete="off"
            @focus="onSegmentFocus(index)"
            @blur="onSegmentBlur"
            @keydown="onKeyDown($event, index)"
            @input="onInput($event, index)"
            @click="onSegmentClick(index)"
          />
        </template>
        <button v-if="modelValue && clearable" type="button" class="datetime-clear" :aria-label="'Limpiar ' + (label || 'fecha/hora')" @click="clearValue">
          <svg width="12" height="12" viewBox="0 0 14 14" fill="none" aria-hidden="true">
            <path d="M1 1l12 12M13 1L1 13" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </button>
      </div>

      <span :id="inputId + '-hint'" class="datetime-hint">Formato: {{ format }}</span>
      <span v-if="errorMessage" :id="inputId + '-error'" class="datetime-error" role="alert" aria-live="polite">{{ errorMessage }}</span>
      <span :id="inputId + '-correction'" class="sr-only" aria-live="assertive" aria-atomic="true">{{ correctionMessage }}</span>
    </fieldset>
  `,

  props: {
    modelValue: { type: String, default: null },
    format: { type: String, required: true },
    label: { type: String, default: '' },
    required: { type: Boolean, default: false },
    clearable: { type: Boolean, default: true },
    id: { type: String, default: null },
    isoOutput: { type: Boolean, default: true }
  },

  emits: ['update:modelValue', 'change'],

  data() {
    return {
      inputId: this.id || `datetime-input-${++_idCounter}`,
      isFocused: false,
      errorMessage: '',
      correctionMessage: '',
      values: { DD: '', MM: '', YYYY: '', YY: '', hh: '', mm: '', ss: '' },
      segmentRefs: []
    }
  },

  computed: {
    inputType() {
      const hasDate = /DD|MM|YYYY|YY/.test(this.format)
      const hasTime = /hh|mm|ss/.test(this.format)
      if (hasDate && hasTime) return 'datetime'
      if (hasDate) return 'date'
      return 'time'
    },

    segments() {
      const tokens = ['YYYY', 'YY', 'DD', 'MM', 'hh', 'mm', 'ss']
      const tokenInfo = {
        DD:   { length: 2, min: 1,  max: 31,   placeholder: 'DD',   ariaLabel: 'Día' },
        MM:   { length: 2, min: 1,  max: 12,   placeholder: 'MM',   ariaLabel: 'Mes' },
        YYYY: { length: 4, min: 1,  max: 9999, placeholder: 'AAAA', ariaLabel: 'Año' },
        YY:   { length: 2, min: 0,  max: 99,   placeholder: 'AA',   ariaLabel: 'Año corto' },
        hh:   { length: 2, min: 0,  max: 23,   placeholder: 'hh',   ariaLabel: 'Hora' },
        mm:   { length: 2, min: 0,  max: 59,   placeholder: 'mm',   ariaLabel: 'Minutos' },
        ss:   { length: 2, min: 0,  max: 59,   placeholder: 'ss',   ariaLabel: 'Segundos' }
      }
      const result = []
      let remaining = this.format
      while (remaining.length > 0) {
        let matched = false
        for (const token of tokens) {
          if (remaining.startsWith(token)) {
            const info = tokenInfo[token]
            result.push({ type: 'field', token, ...info, display: this.values[token], numericValue: this.values[token] ? parseInt(this.values[token]) : null })
            remaining = remaining.slice(token.length)
            matched = true
            break
          }
        }
        if (!matched) {
          result.push({ type: 'separator', char: remaining[0] })
          remaining = remaining.slice(1)
        }
      }
      return result
    },

    fieldIndices() {
      return this.segments.map((s, i) => s.type === 'field' ? i : -1).filter(i => i >= 0)
    },

    descriptionIds() {
      const ids = [`${this.inputId}-hint`]
      if (this.errorMessage) ids.push(`${this.inputId}-error`)
      return ids.join(' ')
    }
  },

  watch: {
    modelValue: { immediate: true, handler(val) { this.parseModelValue(val) } }
  },

  methods: {
    parseModelValue(val) {
      if (!val) { this.values = { DD: '', MM: '', YYYY: '', YY: '', hh: '', mm: '', ss: '' }; return }

      if (!this.isoOutput) {
        // Extraer valores posicionalmente según el format
        const tokens = ['YYYY', 'YY', 'DD', 'MM', 'hh', 'mm', 'ss']
        let remaining = this.format
        let valRemaining = val
        while (remaining.length > 0) {
          let matched = false
          for (const token of tokens) {
            if (remaining.startsWith(token)) {
              const len = token === 'YYYY' ? 4 : 2
              this.values[token] = valRemaining.slice(0, len)
              remaining = remaining.slice(token.length)
              valRemaining = valRemaining.slice(len)
              matched = true
              break
            }
          }
          if (!matched) { remaining = remaining.slice(1); valRemaining = valRemaining.slice(1) }
        }
        return
      }
      
      let day='', month='', year='', hour='', minute='', second=''
      if (this.inputType === 'date' || this.inputType === 'datetime') {
        const datePart = val.includes('T') ? val.split('T')[0] : val.split(' ')[0]
        const parts = datePart.split('-')
        if (parts.length >= 3) { year = parts[0]; month = parts[1]; day = parts[2] }
      }
      if (this.inputType === 'time' || this.inputType === 'datetime') {
        const timePart = val.includes('T') ? val.split('T')[1] : (val.includes(' ') ? val.split(' ')[1] : val)
        const parts = timePart.split(':')
        if (parts.length >= 1) hour = parts[0]
        if (parts.length >= 2) minute = parts[1]
        if (parts.length >= 3) second = parts[2]
      }
      this.values = { DD: day, MM: month, YYYY: year, YY: year ? year.slice(-2) : '', hh: hour, mm: minute, ss: second }
    },

    emitValue() {
      const requiredTokens = this.segments.filter(s => s.type === 'field').map(s => s.token)
      const allFilled = requiredTokens.every(token => {
        const val = this.values[token]
        return val !== '' && val.length === this.getTokenLength(token)
      })
      if (!allFilled) { this.$emit('update:modelValue', null); return }
      if (this.inputType !== 'time') { this.validateDateExistence() }
      try {
        let isoValue
        if (!this.isoOutput) {
          // Reconstruir el valor respetando el format visual
          isoValue = this.format
          isoValue = isoValue.replace('DD',   this.values.DD.padStart(2,'0'))
          isoValue = isoValue.replace('MM',   this.values.MM.padStart(2,'0'))
          isoValue = isoValue.replace('YYYY', this.values.YYYY.padStart(4,'0'))
          isoValue = isoValue.replace('YY',   this.values.YY.padStart(2,'0'))
          isoValue = isoValue.replace('hh',   this.values.hh.padStart(2,'0'))
          isoValue = isoValue.replace('mm',   this.values.mm.padStart(2,'0'))
          isoValue = isoValue.replace('ss',   this.values.ss.padStart(2,'0'))
        } else if (this.inputType === 'date') {
          const y = this.values.YYYY || `20${this.values.YY}`
          isoValue = `${y}-${this.values.MM.padStart(2,'0')}-${this.values.DD.padStart(2,'0')}`
        } else if (this.inputType === 'time') {
          const h = (this.values.hh||'00').padStart(2,'0')
          const m = (this.values.mm||'00').padStart(2,'0')
          const s = (this.values.ss||'00').padStart(2,'0')
          isoValue = requiredTokens.includes('ss') ? `${h}:${m}:${s}` : `${h}:${m}`
        } else {
          const y = this.values.YYYY || `20${this.values.YY}`
          const date = `${y}-${this.values.MM.padStart(2,'0')}-${this.values.DD.padStart(2,'0')}`
          const h = (this.values.hh||'00').padStart(2,'0')
          const m = (this.values.mm||'00').padStart(2,'0')
          const s = requiredTokens.includes('ss') ? `:${(this.values.ss||'00').padStart(2,'0')}` : ''
          isoValue = `${date}T${h}:${m}${s}`
        }
        this.errorMessage = ''
        this.$emit('update:modelValue', isoValue)
        this.$emit('change', isoValue)
      } catch(e) { this.errorMessage = 'Fecha/hora inválida' }

    },

    validateDateExistence() {
      const dd = parseInt(this.values.DD)
      const mm = parseInt(this.values.MM)
      const yyyy = this.values.YYYY
        ? parseInt(this.values.YYYY)
        : (this.values.YY ? 2000 + parseInt(this.values.YY) : null)
      if (!dd || !mm) return
      const year = yyyy || 2000
      const daysInMonth = new Date(year, mm, 0).getDate()
      if (dd > daysInMonth) {
        const originalDay = dd
        this.values.DD = String(daysInMonth).padStart(2, '0')
        const monthNames = ['enero','febrero','marzo','abril','mayo','junio',
          'julio','agosto','septiembre','octubre','noviembre','diciembre']
        const monthName = monthNames[mm - 1] || `mes ${mm}`
        const yearText = yyyy ? ` de ${yyyy}` : ''
        this.correctionMessage = ''
        this.$nextTick(() => {
          this.correctionMessage =
            `Anteriormente has escrito una fecha inválida. El día ${originalDay} no existe en ${monthName} ${yearText}. ` +
            `Lo he corregido al último día válido: ${daysInMonth}`
        })
      }
    },

    getTokenLength(token) { return token === 'YYYY' ? 4 : 2 },

    onSegmentFocus(index) {
      this.isFocused = true
      this.$nextTick(() => { const el = this.segmentRefs[index]; if (el) el.select() })
    },

    onSegmentBlur() {
      setTimeout(() => {
        const focused = document.activeElement
        const isInside = this.segmentRefs.some(ref => ref === focused)
        if (!isInside) { this.isFocused = false; this.padSegments(); this.emitValue() }
      }, 100)
    },

    onSegmentClick(index) {
      this.$nextTick(() => { const el = this.segmentRefs[index]; if (el) el.select() })
    },

    onKeyDown(event, index) {
      const segment = this.segments[index]
      if (segment.type !== 'field') return
      const { key } = event
      if (key === 'ArrowUp' || key === 'ArrowDown') { event.preventDefault(); this.adjustValue(index, key === 'ArrowUp' ? 1 : -1); return }
      if (key === 'ArrowRight') { event.preventDefault(); this.focusNextField(index); return }
      if (key === 'ArrowLeft') { event.preventDefault(); this.focusPrevField(index); return }
      if (key === 'Tab') return
      if (key === 'Backspace' || key === 'Delete') { event.preventDefault(); this.values[segment.token] = ''; return }
      if (!/^\d$/.test(key)) { event.preventDefault(); return }
    },

    onInput(event, index) {
      const segment = this.segments[index]
      if (segment.type !== 'field') return
      const rawValue = event.target.value.replace(/\D/g, '')
      const token = segment.token
      const maxLen = this.getTokenLength(token)
      const truncated = rawValue.slice(0, maxLen)
      this.values[token] = truncated
      if (truncated.length === maxLen) {
        this.values[token] = this.validateAndClamp(token, truncated)
        if ((token === 'DD' || token === 'MM') && this.inputType !== 'time') {
          const hasDd = this.values.DD.length === 2
          const hasMm = this.values.MM.length === 2
          if (hasDd && hasMm) { this.validateDateExistence() }
        }
        this.focusNextField(index)
        this.emitValue()
      }
    },

    validateAndClamp(token, value) {
      const segment = this.segments.find(s => s.type === 'field' && s.token === token)
      if (!segment) return value
      const num = parseInt(value)
      if (isNaN(num)) return value
      const clamped = Math.min(Math.max(num, segment.min), segment.max)
      const result = String(clamped).padStart(this.getTokenLength(token), '0')
      if (clamped !== num) {
        this.announceCorrection(segment.ariaLabel, num, clamped, segment.min, segment.max)
      }
      return result
    },

    announceCorrection(fieldLabel, entered, corrected, min, max) {
      this.correctionMessage = ''
      this.$nextTick(() => {
        this.correctionMessage =
          `Anteriormente has escrito un valor fuera de rango en el campo ${fieldLabel}. Escribiste ${entered} pero el rango permitido es de ${min} a ${max}. Lo he corregido a ${corrected}.`
      })
    },

    adjustValue(index, delta) {
      const segment = this.segments[index]
      if (segment.type !== 'field') return
      const token = segment.token
      const current = parseInt(this.values[token]) || 0
      const clamped = Math.min(Math.max(current + delta, segment.min), segment.max)
      this.values[token] = String(clamped).padStart(this.getTokenLength(token), '0')
      this.emitValue()
    },

    padSegments() {
      for (const segment of this.segments) {
        if (segment.type !== 'field') continue
        const val = this.values[segment.token]
        if (val && val.length > 0 && val.length < this.getTokenLength(segment.token)) {
          this.values[segment.token] = this.validateAndClamp(segment.token, val)
        }
      }
    },

    focusNextField(currentIndex) {
      const nextIdx = this.fieldIndices.find(i => i > currentIndex)
      if (nextIdx !== undefined && this.segmentRefs[nextIdx]) { this.segmentRefs[nextIdx].focus(); this.segmentRefs[nextIdx].select() }
    },

    focusPrevField(currentIndex) {
      const prevIndices = this.fieldIndices.filter(i => i < currentIndex)
      if (prevIndices.length > 0) {
        const prevIdx = prevIndices[prevIndices.length - 1]
        if (this.segmentRefs[prevIdx]) { this.segmentRefs[prevIdx].focus(); this.segmentRefs[prevIdx].select() }
      }
    },

    clearValue() {
      this.values = { DD: '', MM: '', YYYY: '', YY: '', hh: '', mm: '', ss: '' }
      this.errorMessage = ''
      this.$emit('update:modelValue', null)
      this.$emit('change', null)
      this.$nextTick(() => {
        const firstFieldIdx = this.fieldIndices[0]
        if (firstFieldIdx !== undefined && this.segmentRefs[firstFieldIdx]) this.segmentRefs[firstFieldIdx].focus()
      })
    }
  }
}