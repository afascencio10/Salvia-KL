/*
  Botón "Volver" transversal.

  Antes cada módulo navegaba a un destino fijo (NAVIGATION_RULES / breadcrumb
  quemado), por lo que el usuario terminaba en el inicio del módulo en lugar de
  la pantalla inmediatamente anterior.

  salviaGoBack(fallback) usa el historial real del navegador y sólo cae al
  destino fijo cuando no hay una pantalla anterior utilizable (entrada directa
  por URL, pestaña nueva, o venir del login).
*/
(function (root) {

  var EXCLUDED = /^\/(security\/|static\/landing)/;

  // ¿Es seguro hacer history.back()? Puro, para poder probarlo.
  function canGoBack(referrer, currentUrl, historyLength) {
    if (!referrer || historyLength <= 1) return false;
    try {
      var ref = new URL(referrer);
      var cur = new URL(currentUrl);
      if (ref.origin !== cur.origin) return false;      // salió de la app
      if (EXCLUDED.test(ref.pathname)) return false;    // login / landing
      if (ref.pathname === cur.pathname) return false;  // misma pantalla (recarga, guardado)
      return true;
    } catch (e) {
      return false;
    }
  }

  root.salviaCanGoBack = canGoBack;

  root.salviaGoBack = function (fallback) {
    if (canGoBack(document.referrer, location.href, history.length)) {
      history.back();
      return;
    }
    location.assign(fallback || '/');
  };

  // Autocomprobación: node src/frontend/js/nav-back.js
  if (typeof module !== 'undefined' && require.main === module) {
    var assert = require('assert');
    var here = 'https://app/salvia/casos/ABC/detalle';
    assert.strictEqual(canGoBack('https://app/salvia/lista-casos', here, 3), true);
    assert.strictEqual(canGoBack('', here, 3), false, 'sin referrer');
    assert.strictEqual(canGoBack('https://app/salvia/lista-casos', here, 1), false, 'pestaña nueva');
    assert.strictEqual(canGoBack('https://otro/x', here, 3), false, 'otro origen');
    assert.strictEqual(canGoBack('https://app/static/landing.html', here, 3), false, 'landing');
    assert.strictEqual(canGoBack('https://app/security/login', here, 3), false, 'login');
    assert.strictEqual(canGoBack(here, here, 3), false, 'misma pantalla');
    console.log('nav-back OK');
  }

})(typeof window !== 'undefined' ? window : globalThis);
