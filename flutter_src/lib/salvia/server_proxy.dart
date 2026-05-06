import 'package:demo_flut/salvia/victim_contact.dart';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:io';
import 'dart:typed_data';

import 'package:http/io_client.dart';

class ServerProxy {
  static String cookies = "";
  static final Map<String, String> _cookieJar = {};
  static const bool _debugCaptcha = true;
  static const String baseUrl = String.fromEnvironment(
    'SALVIA_BASE_URL',
    defaultValue: 'https://pruebassalvia.minigualdadyequidad.gov.co',
  );

  static String get _normalizedBaseUrl => baseUrl.endsWith('/')
      ? baseUrl.substring(0, baseUrl.length - 1)
      : baseUrl;

  static Uri _uri(String path) {
    final normalizedPath = path.startsWith('/') ? path : '/$path';
    return Uri.parse(_normalizedBaseUrl).resolve(normalizedPath);
  }

  static String captchaImageUrl(
    String captchaId, {
    int cacheBust = 0,
    bool legacy = false,
  }) {
    final query = cacheBust > 0 ? '?v=$cacheBust' : '';
    final suffix = legacy ? '' : '/img';
    return '${_uri('/seguridad/login/captcha/$captchaId$suffix')}$query';
  }

  static String captchaAudioUrl(String captchaId, {int cacheBust = 0}) {
    final query = cacheBust > 0 ? '?v=$cacheBust' : '';
    return '${_uri('/seguridad/login/captcha/$captchaId/audio')}$query';
  }

  static void _captureCookies(String? setCookieHeader) {
    if (setCookieHeader == null || setCookieHeader.trim().isEmpty) {
      return;
    }

    final matches =
        RegExp(r'(^|,\s*)([^=;,\s]+)=([^;,\r\n]+)').allMatches(setCookieHeader);
    for (final m in matches) {
      final name = m.group(2);
      final value = m.group(3);
      if (name == null || value == null) continue;

      final lower = name.toLowerCase();
      if (lower == 'path' ||
          lower == 'expires' ||
          lower == 'max-age' ||
          lower == 'domain' ||
          lower == 'secure' ||
          lower == 'httponly' ||
          lower == 'samesite') {
        continue;
      }
      _cookieJar[name] = value;
    }

    if (_cookieJar.containsKey('default_session')) {
      cookies = 'default_session=${_cookieJar['default_session']}';
      return;
    }

    cookies = _cookieJar.entries.map((e) => '${e.key}=${e.value}').join('; ');
  }

  static String _cookieHeader() {
    final fromJar =
        _cookieJar.entries.map((e) => '${e.key}=${e.value}').join('; ');
    if (fromJar.trim().isNotEmpty) return fromJar;
    return cookies.trim();
  }

  static void _logCaptcha(String message) {
    if (!_debugCaptcha) return;
    // ignore: avoid_print
    print('[captcha] $message');
  }

  static String _decodeUtf8Body(http.Response response) {
    return utf8.decode(response.bodyBytes, allowMalformed: true);
  }

  static HttpClient getUnsafeHttpClient() {
    if (kIsWeb) {
      throw UnsupportedError('HttpClient is not available on Flutter Web');
    }
    HttpClient httpClient = HttpClient()
      ..badCertificateCallback =
          ((X509Certificate cert, String host, int port) => true);
    return httpClient;
  }

  static http.Client _buildClient() {
    if (kIsWeb) {
      return http.Client();
    }
    final ioc = HttpClient()
      ..badCertificateCallback =
          (X509Certificate cert, String host, int port) => true;
    return IOClient(ioc);
  }

  Future<RespuestaServidor> login(String usuario, String password) async {
    RespuestaServidor res =
        RespuestaServidor(success: false, mensaje: "Sin conexión");
    Uri url = _uri('/seguridad/login');
    Map<String, String> headers = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };
    String body = jsonEncode({
      'field1': 'valor1',
      'field2': 'valor2',
      'field3': 'valor3',
    });

    try {
      final client = _buildClient();
      final response = await client.post(url, headers: headers, body: body);
      client.close();
      _captureCookies(response.headers['set-cookie']);
      res = RespuestaServidor.fromResponse(response);
      return res;
    } catch (e) {
      return res;
    }
  }

  Future<RespuestaServidor> setReport(VictimContact victimContact) async {
    RespuestaServidor res =
        RespuestaServidor(success: false, mensaje: "Sin conexión");
    Uri url = _uri('/salvia/public');
    final cookieValue = _cookieHeader();
    Map<String, String> headers = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };
    if (cookieValue.isNotEmpty) {
      headers['Cookie'] = cookieValue;
    }
    String body = jsonEncode({'victimContact': victimContact.toJson()});

    try {
      final client = _buildClient();
      final response = await client.post(url, headers: headers, body: body);
      client.close();
      res = RespuestaServidor.fromResponse(response);
      _captureCookies(response.headers['set-cookie']);
      return res;
    } catch (e) {
      res.mensaje = e.toString();
      return res;
    }
  }

  Future<RespuestaServidor> submitPrimerContacto(
    Map<String, dynamic> payload,
  ) async {
    RespuestaServidor res =
        RespuestaServidor(success: false, mensaje: "Sin conexión");
    Uri url = _uri('/public/primer_contacto/nuevo');
    final cookieValue = _cookieHeader();
    Map<String, String> headers = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };
    if (cookieValue.isNotEmpty) {
      headers['Cookie'] = cookieValue;
    }
    String body = jsonEncode(payload);

    try {
      final client = _buildClient();
      final response = await client.post(url, headers: headers, body: body);
      client.close();
      res = RespuestaServidor.fromResponse(response);

      _captureCookies(response.headers['set-cookie']);

      return res;
    } catch (e) {
      res.mensaje = e.toString();
      return res;
    }
  }

  Future<String?> fetchCaptchaIdFromPublicForm() async {
    final config = await fetchPublicFormConfig();
    if (config != null && config.captchaId.trim().isNotEmpty) {
      return config.captchaId.trim();
    }
    return _fetchCaptchaIdFromPublicFormHtmlFallback();
  }

  Future<PublicFormConfig?> fetchPublicFormConfig() async {
    final url = _uri('/salvia/public/nuevo');
    _logCaptcha('fetchPublicFormConfig base=$_normalizedBaseUrl');

    final client = _buildClient();
    try {
      final responseJson = await client.get(
        url,
        headers: const {
          'Accept': 'application/json',
        },
      );
      _captureCookies(responseJson.headers['set-cookie']);
      final responseJsonText = _decodeUtf8Body(responseJson);
      _logCaptcha(
        'GET $url (json) status=${responseJson.statusCode} len=${responseJsonText.length}',
      );

      if (responseJson.statusCode != 200) {
        return null;
      }

      final decoded = jsonDecode(responseJsonText);
      if (decoded is! Map) return null;
      final fields = decoded['fields'];
      if (fields is! Map) return null;

      final fieldMap = Map<String, dynamic>.from(fields);
      final captchaId = (fieldMap['captchaID'] ?? '').toString().trim();

      final yesNoRaw = fieldMap['yesNo'] ?? fieldMap['yes_no'];
      final reportTypeRaw =
          fieldMap['reportType'] ?? fieldMap['victim_case_form2_report_type'];
      final reportTypeDetailsRaw = fieldMap['reportTypeDetails'] ??
          fieldMap['victim_case_form2_report_type_details'];
      final adjustmentsGBVRaw = fieldMap['adjustmentsGBV'] ??
          fieldMap['victim_case_form2_adjustments_gbv'];

      final yesNo = _asListOfMap(yesNoRaw);
      final reportType = _asListOfMap(reportTypeRaw);
      final reportTypeDetails = _asListOfMap(reportTypeDetailsRaw);
      final adjustmentsGBV = _asListOfMap(adjustmentsGBVRaw);

      _logCaptcha(
        'public form config captchaID="$captchaId" yesNo=${yesNo.length} reportType=${reportType.length} reportTypeDetails=${reportTypeDetails.length} adjustmentsGBV=${adjustmentsGBV.length}',
      );

      return PublicFormConfig(
        captchaId: captchaId,
        yesNo: yesNo,
        reportType: reportType,
        reportTypeDetails: reportTypeDetails,
        adjustmentsGBV: adjustmentsGBV,
      );
    } catch (e) {
      _logCaptcha('fetchPublicFormConfig exception: $e');
      return null;
    } finally {
      client.close();
    }
  }

  Future<String?> _fetchCaptchaIdFromPublicFormHtmlFallback() async {
    final url = _uri('/salvia/public/nuevo');
    final client = _buildClient();
    try {
      final responseHtml = await client.get(
        url,
        headers: const {
          'Accept': 'text/html',
        },
      );
      _captureCookies(responseHtml.headers['set-cookie']);
      final responseHtmlText = _decodeUtf8Body(responseHtml);
      _logCaptcha(
        'GET $url (html) status=${responseHtml.statusCode} len=${responseHtmlText.length}',
      );
      if (responseHtml.statusCode != 200) {
        return null;
      }

      final html = responseHtmlText;
      final match = RegExp(r'captchaID\s*:\s*"([^"]+)"').firstMatch(html);
      final htmlId = match?.group(1);
      _logCaptcha('captchaID(html)="${htmlId ?? ''}"');
      return htmlId;
    } catch (e) {
      _logCaptcha('html fallback exception: $e');
      return null;
    } finally {
      client.close();
    }
  }

  List<Map<String, dynamic>> _asListOfMap(dynamic raw) {
    if (raw is! List) return const [];
    final out = <Map<String, dynamic>>[];
    for (final item in raw) {
      if (item is Map<String, dynamic>) {
        out.add(item);
      } else if (item is Map) {
        out.add(Map<String, dynamic>.from(item));
      }
    }
    return out;
  }

  Future<Uint8List?> fetchCaptchaImageBytes(
    String captchaId, {
    int cacheBust = 0,
    int retries = 1,
  }) async {
    final urls = [
      Uri.parse(captchaImageUrl(captchaId, cacheBust: cacheBust)),
      Uri.parse(captchaImageUrl(captchaId, cacheBust: cacheBust, legacy: true)),
    ];
    final cookieValue = _cookieHeader();
    final headers = {
      'Accept': 'image/png,image/*;q=0.8,*/*;q=0.5',
      'Referer': '$_normalizedBaseUrl/salvia/public/nuevo',
      if (cookieValue.isNotEmpty) 'Cookie': cookieValue,
    };

    try {
      final client = _buildClient();
      for (final url in urls) {
        final response = await client.get(url, headers: headers);
        _captureCookies(response.headers['set-cookie']);
        _logCaptcha(
          'GET $url status=${response.statusCode} bytes=${response.bodyBytes.length}',
        );
        if (response.statusCode == 200 && response.bodyBytes.isNotEmpty) {
          client.close();
          return response.bodyBytes;
        }
      }
      client.close();
      if (retries > 0) {
        _logCaptcha('captcha empty, retrying ($retries left)');
        return fetchCaptchaImageBytes(
          captchaId,
          cacheBust: DateTime.now().millisecondsSinceEpoch,
          retries: retries - 1,
        );
      }
      _logCaptcha('captcha image still empty after retries');
      return null;
    } catch (e) {
      _logCaptcha('fetchCaptchaImageBytes exception: $e');
      return null;
    }
  }

  Future<Uint8List?> fetchCaptchaAudioBytes(
    String captchaId, {
    int cacheBust = 0,
    int retries = 1,
  }) async {
    final query = cacheBust > 0 ? '?v=$cacheBust' : '';
    final url =
        Uri.parse('${_uri('/seguridad/login/captcha/$captchaId/audio')}$query');
    final cookieValue = _cookieHeader();
    final headers = {
      'Accept': 'audio/wav,audio/*;q=0.8,*/*;q=0.5',
      'Referer': '$_normalizedBaseUrl/salvia/public/nuevo',
      if (cookieValue.isNotEmpty) 'Cookie': cookieValue,
    };

    try {
      final client = _buildClient();
      final response = await client.get(url, headers: headers);
      _captureCookies(response.headers['set-cookie']);
      _logCaptcha(
        'GET $url status=${response.statusCode} bytes=${response.bodyBytes.length} type=${response.headers['content-type'] ?? ''}',
      );
      client.close();

      if (response.statusCode == 200 && response.bodyBytes.isNotEmpty) {
        return response.bodyBytes;
      }

      if (retries > 0) {
        _logCaptcha('captcha audio empty, retrying ($retries left)');
        return fetchCaptchaAudioBytes(
          captchaId,
          cacheBust: DateTime.now().millisecondsSinceEpoch,
          retries: retries - 1,
        );
      }
      _logCaptcha('captcha audio still empty after retries');
      return null;
    } catch (e) {
      _logCaptcha('fetchCaptchaAudioBytes exception: $e');
      return null;
    }
  }
}

class PublicFormConfig {
  final String captchaId;
  final List<Map<String, dynamic>> yesNo;
  final List<Map<String, dynamic>> reportType;
  final List<Map<String, dynamic>> reportTypeDetails;
  final List<Map<String, dynamic>> adjustmentsGBV;

  const PublicFormConfig({
    required this.captchaId,
    required this.yesNo,
    required this.reportType,
    required this.reportTypeDetails,
    required this.adjustmentsGBV,
  });
}

class RespuestaServidor {
  // Atributos públicos
  bool success;
  int statusCode;
  String mensaje;
  Map<String, dynamic> data;

  // Constructor
  RespuestaServidor(
      {this.success = false,
      this.statusCode = 0,
      this.mensaje = '',
      this.data = const {}});

  // Constructor con nombre para parsear desde JSON
  RespuestaServidor.fromJsonObj(Map<String, dynamic> json)
      : success = json['status'] ?? false,
        statusCode = 0,
        mensaje = json['msg'] ?? '',
        data = const {};

  RespuestaServidor.fromResponse(http.Response res)
      : success = res.statusCode == 200,
        statusCode = res.statusCode,
        mensaje = _decodeResponseBody(res),
        data = _tryDecode(_decodeResponseBody(res));

  static String _decodeResponseBody(http.Response res) {
    return utf8.decode(res.bodyBytes, allowMalformed: true);
  }

  static Map<String, dynamic> _tryDecode(String raw) {
    try {
      final decoded = jsonDecode(raw);
      if (decoded is Map<String, dynamic>) {
        return decoded;
      }
      return {'raw': decoded};
    } catch (_) {
      return {};
    }
  }
}
