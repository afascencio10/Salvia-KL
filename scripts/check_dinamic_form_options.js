/* Check de las opciones dinámicas de dinamic-form (state_options_path).
 *
 * Corre con:  node scripts/check_dinamic_form_options.js
 *
 * Cubre las dos piezas de las que dependen Ciudad, Municipio y Competencias del
 * repeater de barreras: que el path con {_entryIndex} resuelva contra el formState
 * del padre, y que una lista de opciones vacía NO borre la respuesta guardada
 * (era pérdida de datos silenciosa al reabrir un seguimiento).
 *
 * Carga el componente real en un sandbox con un `app` de mentira, así que prueba
 * el código que se sirve, no una copia.
 */

const assert = require('assert');
const fs = require('fs');
const path = require('path');
const vm = require('vm');

const componentPath = process.argv[2]
	|| path.join(__dirname, '..', 'src', 'frontend', 'js', 'components', 'dinamic-form.js');
const source = fs.readFileSync(componentPath, 'utf8');

const sandbox = { console, definition: null };
sandbox.app = { component: (_name, definition) => { sandbox.definition = definition; } };
// El componente inyecta sus estilos al cargar; devolver algo truthy lo hace salir
// temprano sin necesitar un DOM.
sandbox.document = { getElementById: () => ({}) };
vm.runInNewContext(source + '\n;globalThis.__resolve = resolveQuestionOptions;', sandbox);

const resolveQuestionOptions = sandbox.__resolve;
const { _sanitizeOptionAnswer } = sandbox.definition.methods;
assert.ok(typeof resolveQuestionOptions === 'function', 'resolveQuestionOptions no se pudo cargar');
assert.ok(typeof _sanitizeOptionAnswer === 'function', '_sanitizeOptionAnswer no se pudo cargar');

const competencias = {
	id: 'q-competencias',
	questionTypeId: 'multiple',
	stateOptionsPath: 'newBarriers.{_entryIndex}.competencias',
};
const formState = {
	newBarriers: [
		{ cities: [], towns: [], competencias: [] },
		{ cities: [], towns: [], competencias: [{ label: 'Verificar afiliación', value: 'ips_aaa' }, { label: 'Remitir a Fiscalía', value: 'ips_bbb' }] },
	],
};

const checks = [
	['el path con {_entryIndex} resuelve la entrada correcta', () => {
		assert.strictEqual(resolveQuestionOptions(competencias, formState, 1).length, 2);
		assert.strictEqual(resolveQuestionOptions(competencias, formState, 0).length, 0);
	}],
	['sin opciones cargadas la respuesta guardada se conserva', () => {
		// entryIndex 0 todavía no tiene competencias: el fetch no ha resuelto.
		assert.strictEqual(
			_sanitizeOptionAnswer(competencias, 'ips_aaa,ips_bbb', formState, 0),
			'ips_aaa,ips_bbb',
		);
	}],
	['con opciones cargadas se filtran las que ya no existen', () => {
		assert.strictEqual(
			_sanitizeOptionAnswer(competencias, 'ips_aaa,ips_borrada', formState, 1),
			'ips_aaa',
		);
	}],
	['single con valor inválido se limpia', () => {
		const ciudad = { id: 'q-ciudad', questionTypeId: 'dropdown', stateOptionsPath: 'newBarriers.{_entryIndex}.competencias' };
		assert.strictEqual(_sanitizeOptionAnswer(ciudad, 'ips_aaa', formState, 1), 'ips_aaa');
		assert.strictEqual(_sanitizeOptionAnswer(ciudad, 'no_existe', formState, 1), '');
	}],
	['una pregunta sin state_options_path no se toca', () => {
		const texto = { id: 'q-texto', questionTypeId: 'text' };
		assert.strictEqual(_sanitizeOptionAnswer(texto, 'lo que sea', formState, 1), 'lo que sea');
	}],
];

let failed = 0;
for (const [name, check] of checks) {
	try {
		check();
		console.log(`  ok   ${name}`);
	} catch (error) {
		failed++;
		console.error(`  FALLA ${name}\n        ${error.message}`);
	}
}
console.log(failed === 0 ? `\n${checks.length} checks OK` : `\n${failed} de ${checks.length} fallaron`);
process.exit(failed === 0 ? 0 : 1);
