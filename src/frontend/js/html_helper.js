var buildInput = function (inputId, errorSpanId, inputLabel, inputPlaceHolder, entityName, entityAttr) {
  return `
  <div class="form-group">
    <label for="${inputId}">${inputLabel}</label>

    <template v-if="errors.${entityName}&&errors.${entityName}.${entityAttr}">
      <input v-model="data.${entityName}.${entityAttr}" type="text" class="form-control is-invalid" id="${inputId}" placeholder="${inputPlaceHolder}">
      <span id="${errorSpanId}" class="error invalid-feedback">\${errors.${entityName}.${entityAttr}}</span>
    </template>
    <input v-else v-model="data.${entityName}.${entityAttr}" type="text" class="form-control" id="${inputId}" placeholder="${inputPlaceHolder}">
  </div>`;
};

var buildSelect = function (selectId, errorSpanId, selectLabel, selectPlaceHolder, entityName, entityAttr, enumsAttr) {
  return `
  <div class="form-group">
    <label for="${selectId}">${selectLabel}</label>

    <template v-if="errors.${entityName}&&errors.${entityName}.${entityAttr}">
      <select v-model="data.${entityName}.${entityAttr}" class="custom-select form-control-border border-width-2 is-invalid" id="${selectId}">
        <option v-for="(label, key) in enums.${enumsAttr}" :value="key">\${label}</option>
      </select>
      <span id="${errorSpanId}" class="error invalid-feedback">\${errors.${entityName}.${entityAttr}}</span>
    </template>
    <select v-else v-model="data.${entityName}.${entityAttr}" class="custom-select form-control-border border-width-2" id="${selectId}">
      <option v-for="(label, key) in enums.${enumsAttr}" :value="key">\${label}</option>
    </select>
  </div>`;
};