import 'dart:async';
import 'dart:convert';

import 'package:audioplayers/audioplayers.dart';
import 'package:demo_flut/salvia/server_proxy.dart';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:geolocator/geolocator.dart';

import '../global/my_colors.dart';
import '../global/my_screens.dart';
import '../widgets/elements/espacio.dart';

class ReportarScreen extends StatefulWidget {
  const ReportarScreen({super.key});

  @override
  State<ReportarScreen> createState() => _ReportarScreenState();
}

class _ReportarScreenState extends State<ReportarScreen> {
  static const bool _skipCaptcha = bool.fromEnvironment(
    'SALVIA_SKIP_CAPTCHA',
    defaultValue: false,
  );
  static const double _fallbackLatitude = 4.60971;
  static const double _fallbackLongitude = -74.08175;
  static const double _sectionHeaderFontSize = 20.0;
  static const double _fieldLabelFontSize = 13.0;
  static const int _headerTopFlex = 42;
  static const int _headerCenefaFlex = 6;
  static const int _headerBottomFlex = 46;
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  final ScrollController _scrollController = ScrollController();

  final TextEditingController _reporterNamesCtrl = TextEditingController();
  final TextEditingController _reporterPhoneCtrl = TextEditingController();
  final TextEditingController _victimNamesCtrl = TextEditingController();
  final TextEditingController _victimLastNamesCtrl = TextEditingController();
  final TextEditingController _victimColPhoneCtrl = TextEditingController();
  final TextEditingController _bestContactTimeCtrl = TextEditingController();
  final TextEditingController _bestContactTimeDisplayCtrl =
      TextEditingController();
  final TextEditingController _factsCtrl = TextEditingController();
  final TextEditingController _captchaCtrl = TextEditingController();
  final GlobalKey _captchaSectionKey = GlobalKey();

  static const List<_EnumOption> _fallbackYesNoOptions = [
    _EnumOption(
      icode: 'cea7f0d1-32f6-4f91-991b-413b2541851c',
      name: 'Sí',
      code: 'y',
    ),
    _EnumOption(
      icode: 'cdc5eb69-fede-4481-b93d-8aeea9b7b505',
      name: 'No',
      code: 'n',
    ),
  ];

  static const List<_EnumOption> _fallbackReportTypeOptions = [
    _EnumOption(
      icode: 'a35e75c0-a62e-4413-a01c-d2bccd184f33',
      name: 'Usted es la víctima de violencia basada en género',
      code: 'v',
    ),
    _EnumOption(
      icode: '42db617d-b843-4180-977d-305bcad484c0',
      name:
          'Usted quiere reportar un hecho de violencia basada en género experimentado por otra persona',
      code: 'r',
    ),
  ];

  static const List<_EnumOption> _fallbackReportTypeDetailsOptions = [
    _EnumOption(
      icode: 'e8f7a9d0-0969-466b-95af-97e18e294e88',
      name: 'Educadora par ASP',
      code: 'ae',
    ),
    _EnumOption(
      icode: '6c63ca19-5b43-4e0f-a7a8-d0cb597ac044',
      name: 'Enlace Territorial Salvia',
      code: 'et',
    ),
    _EnumOption(
      icode: '2b90265d-352b-46f0-b29e-059b517aff69',
      name: 'Persona Servidora MIyE externa a Salvia',
      code: 'ex',
    ),
    _EnumOption(
      icode: 'de9273ce-48e2-4d4e-9267-74645c3d95a1',
      name: 'Otra',
      code: 'ot',
    ),
    _EnumOption(
      icode: 'b1b2ff03-5fff-4885-941c-52885a23fa69',
      name: 'Anónima',
      code: 'an',
    ),
  ];

  static const List<_EnumOption> _fallbackAdjustmentsGbvOptions = [
    _EnumOption(
      icode: '6f133139-ca3c-4321-be35-5ab149bd0201',
      name: 'Intérprete de lengua de señas colombiana',
      code: 'sc',
    ),
    _EnumOption(
      icode: 'a653ecde-b5e2-4d2a-bc95-08c099007e5c',
      name: 'Guía-intérprete para personas sordociegas',
      code: 'gi',
    ),
    _EnumOption(
      icode: 'a2ea378f-cc1d-46b3-b041-46d8724c9fd8',
      name:
          'Intérprete de idiomas / traducción (otra lengua distinta al español)',
      code: 'id',
    ),
    _EnumOption(
      icode: '0f995b93-58a2-4c50-b492-f230f0d91320',
      name:
          'Apoyos para la comunicación (pictogramas, tableros, comunicación aumentativa, etc.)',
      code: 'ap',
    ),
    _EnumOption(
      icode: '5ea5058c-c421-4b5d-ba93-d110c8ac69ce',
      name:
          'Información en formatos accesibles (lectura fácil, macrotipo, audio, braille)',
      code: 'in',
    ),
    _EnumOption(
      icode: '8cbd916c-ec33-4a3e-9874-648d1b2c479c',
      name:
          'Adecuaciones de accesibilidad física (rampas, ascensor, barandas, silla adecuada, etc.)',
      code: 'ad',
    ),
    _EnumOption(
      icode: 'f405a973-9845-4179-9a35-9eacce087e30',
      name:
          'Dispositivos o ayudas técnicas (silla de ruedas, bastón, audífonos, lector de pantalla, etc.)',
      code: 'dt',
    ),
    _EnumOption(
      icode: 'a7a75156-14ba-4f14-b491-a7467ce412ce',
      name:
          'Presencia de una persona de confianza durante la atención y mediación con funcionarios',
      code: 'pr',
    ),
    _EnumOption(
      icode: '58872125-7ec8-45de-bf0c-ed0cd22ed12d',
      name:
          'Transporte seguro o acompañamiento para desplazarse a las instituciones',
      code: 'tr',
    ),
    _EnumOption(
      icode: 'd62602bd-a18f-44b5-bea0-1a0ec61db5da',
      name: 'No informa',
      code: 'ni',
    ),
    _EnumOption(
      icode: '882faac9-bb90-4f17-b977-d01aa9067d22',
      name: 'No requiere',
      code: 'nr',
    ),
    _EnumOption(
      icode: '7831ebda-dcac-40e8-ac44-d915c2ea51e7',
      name: 'No requiere apoyos especiales',
      code: 'ns',
    ),
  ];

  _EnumOption? _authorizationAnswer;
  _EnumOption? _reportType;
  _EnumOption? _reportTypeDetails;
  _EnumOption? _hasCareRole;
  _EnumOption? _victimAwareOfReport;
  _EnumOption? _willReceiveCall;
  List<_EnumOption> _adjustmentsGBVSelected = [];
  List<_EnumOption> _yesNoOptions = List<_EnumOption>.from(
    _fallbackYesNoOptions,
  );
  List<_EnumOption> _reportTypeOptions = List<_EnumOption>.from(
    _fallbackReportTypeOptions,
  );
  List<_EnumOption> _reportTypeDetailsOptions = List<_EnumOption>.from(
    _fallbackReportTypeDetailsOptions,
  );
  List<_EnumOption> _adjustmentsGBVOptions = List<_EnumOption>.from(
    _fallbackAdjustmentsGbvOptions,
  );

  String _captchaId = 'FJwK6ootKKfVFKm71uNM';
  int _captchaCacheBust = 0;
  bool _isRefreshingCaptcha = false;
  bool _isCaptchaLoading = false;
  Uint8List? _captchaBytes;
  final AudioPlayer _captchaAudioPlayer = AudioPlayer();
  StreamSubscription<PlayerState>? _captchaAudioStateSub;
  bool _isCaptchaAudioLoading = false;
  bool _isCaptchaAudioPlaying = false;
  bool _isSubmitting = false;
  double _latitude = _fallbackLatitude;
  double _longitude = _fallbackLongitude;

  static const TextStyle _formPlaceholderStyle = TextStyle(
    color: AppColors.salviaPlaceholderColor,
    fontSize: AppColors.salviaPlaceholderFontSize,
    fontWeight: FontWeight.w400,
  );

  String? _globalError;
  final Map<String, String> _serverErrors = {};

  bool get _isAuthorized => _authorizationAnswer?.code == 'y';
  bool get _showPersonalDataSection => _isAuthorized && _reportType != null;
  bool get _isReporterCase => _reportType?.code == 'r';
  bool get _isVictimCase => _reportType?.code == 'v';
  bool get _hasCareRoleYes => _hasCareRole?.code == 'y';
  bool get _isReportTypeDetailsIdentityCase =>
      const {'ae', 'et', 'ex', 'ot'}.contains(_reportTypeDetails?.code);
  bool get _shouldShowReporterIdentityFields =>
      _isReporterCase && (_hasCareRoleYes || _isReportTypeDetailsIdentityCase);

  final List<TextInputFormatter> _nameInputFormatters = [
    const _NameTextFormatter(maxLength: 32),
  ];

  final List<TextInputFormatter> _fullNameInputFormatters = [
    const _NameTextFormatter(maxLength: 64),
  ];

  final List<TextInputFormatter> _colombiaPhoneInputFormatters = [
    FilteringTextInputFormatter.digitsOnly,
    LengthLimitingTextInputFormatter(10),
  ];

  @override
  void initState() {
    super.initState();
    _captchaAudioStateSub = _captchaAudioPlayer.onPlayerStateChanged.listen((
      state,
    ) {
      if (!mounted) return;
      setState(() {
        _isCaptchaAudioPlaying = state == PlayerState.playing;
      });
    });
    _loadInitialFormAndCaptcha();
    _resolveCurrentLocation();
  }

  @override
  void dispose() {
    _scrollController.dispose();
    _reporterNamesCtrl.dispose();
    _reporterPhoneCtrl.dispose();
    _victimNamesCtrl.dispose();
    _victimLastNamesCtrl.dispose();
    _victimColPhoneCtrl.dispose();
    _bestContactTimeCtrl.dispose();
    _bestContactTimeDisplayCtrl.dispose();
    _factsCtrl.dispose();
    _captchaCtrl.dispose();
    _captchaAudioStateSub?.cancel();
    _captchaAudioPlayer.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        title: const Text(
          'Reportar',
          style: TextStyle(
            color: Color(0xffffffff),
            fontSize: 20,
          ),
        ),
        backgroundColor: AppColors.salviaTransparente,
        elevation: 0,
        automaticallyImplyLeading: true,
        iconTheme: const IconThemeData(color: Colors.white),
      ),
      body: _body(context),
    );
  }

  Widget _body(BuildContext context) {
    final screenHeight = MediaQuery.of(context).size.height;
    final screenWidth = MediaQuery.of(context).size.width;
    final headerHeight = screenHeight * 0.3;
    const totalFlex = _headerTopFlex + _headerCenefaFlex + _headerBottomFlex;
    final panelTop =
        headerHeight * (_headerTopFlex + _headerCenefaFlex) / totalFlex;

    return Stack(
      children: [
        Container(
          decoration: const BoxDecoration(
            color: AppColors.salviaHeaderBottom,
          ),
        ),
        SizedBox(
          width: screenWidth,
          height: headerHeight,
          child: Column(
            children: [
              Expanded(
                flex: _headerTopFlex,
                child: Container(
                  width: double.infinity,
                  color: AppColors.salviaHeaderTopDark,
                ),
              ),
              Expanded(
                flex: _headerCenefaFlex,
                child: Container(
                  width: double.infinity,
                  color: AppColors.salviaHeaderBandSecondary,
                ),
              ),
              Expanded(
                flex: _headerBottomFlex,
                child: Container(
                  width: double.infinity,
                  color: AppColors.salviaHeaderBottom,
                ),
              ),
            ],
          ),
        ),
        SingleChildScrollView(
          controller: _scrollController,
          child: ConstrainedBox(
            constraints: BoxConstraints(
              minHeight: screenHeight,
            ),
            child: Container(
              margin: EdgeInsets.fromLTRB(0, panelTop, 0, 0),
              decoration: BoxDecoration(
                color: const Color(0xffffffff),
                shape: BoxShape.rectangle,
                borderRadius: const BorderRadius.only(
                  topLeft: Radius.circular(30.0),
                  topRight: Radius.circular(30.0),
                ),
                border: Border.all(color: const Color(0x4d9e9e9e), width: 1),
              ),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Form(
                  key: _formKey,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      espacio(40),
                      if (_globalError != null && _globalError!.isNotEmpty)
                        _errorBanner(_globalError!),
                      _sectionHeader(
                        icon: Icons.list_alt,
                        title: 'Autorización de datos personales',
                      ),
                      const SizedBox(height: 14),
                      _dataPolicyAccordion(),
                      const SizedBox(height: 16),
                      _selectField(
                        label:
                            'De acuerdo con lo enunciado, ¿Autoriza el tratamiento de sus datos personal por parte del Ministerio de Igualdad y Equidad? *',
                        placeholder:
                            'Indique si autoriza el tratamiento de datos personales',
                        value: _authorizationAnswer,
                        options: _yesNoOptions,
                        onChanged: (value) {
                          setState(() {
                            _authorizationAnswer = value;
                            _globalError = null;
                            _serverErrors.remove('authorizationAnswer');
                            if (!_isAuthorized) {
                              _adjustmentsGBVSelected = [];
                              _reportType = null;
                              _reportTypeDetails = null;
                              _hasCareRole = null;
                              _victimAwareOfReport = null;
                              _willReceiveCall = null;
                            }
                          });
                        },
                        validator: (value) {
                          if (value == null) {
                            return 'Este campo es obligatorio';
                          }
                          return null;
                        },
                        serverError: _serverErrors['authorizationAnswer'],
                      ),
                      if (_isAuthorized) ...[
                        const SizedBox(height: 16),
                        _multiSelectField(
                          label:
                              '¿Requiere de algún ajuste razonable en el marco de la atención a su caso de VBG y/o VpP? *',
                          placeholder:
                              'Indique si requiere de algún ajuste razonable en el marco de la atención a su caso de VBG y/o VpP',
                          selectedOptions: _adjustmentsGBVSelected,
                          options: _adjustmentsGBVOptions,
                          onChanged: (selected) {
                            setState(() {
                              _adjustmentsGBVSelected = selected;
                              _globalError = null;
                              _serverErrors.remove('adjustmentsGBV');
                            });
                          },
                          validator: (selected) {
                            if (_isAuthorized && selected.isEmpty) {
                              return 'Este campo es obligatorio';
                            }
                            return null;
                          },
                          serverError: _serverErrors['adjustmentsGBV'],
                        ),
                        const SizedBox(height: 16),
                        _selectField(
                          label: 'Tipo de reporte *',
                          placeholder: 'Indique el tipo de reporte',
                          value: _reportType,
                          options: _reportTypeOptions,
                          onChanged: (value) {
                            setState(() {
                              _reportType = value;
                              _globalError = null;
                              _serverErrors.remove('reportType');
                              _reportTypeDetails = null;
                              _hasCareRole = null;
                              _victimAwareOfReport = null;
                              _willReceiveCall = null;
                            });
                          },
                          validator: (value) {
                            if (_isAuthorized && value == null) {
                              return 'Este campo es obligatorio';
                            }
                            return null;
                          },
                          serverError: _serverErrors['reportType'],
                        ),
                      ],
                      if (_showPersonalDataSection) ...[
                        const SizedBox(height: 60),
                        _sectionHeader(
                          icon: Icons.people_alt_rounded,
                          title: 'Datos personales',
                        ),
                        const SizedBox(height: 16),
                        if (_isReporterCase) ...[
                          _selectField(
                            label: 'Usted es una persona *',
                            placeholder: 'Indique si usted es una persona',
                            value: _reportTypeDetails,
                            options: _reportTypeDetailsOptions,
                            onChanged: (value) {
                              setState(() {
                                _reportTypeDetails = value;
                                _globalError = null;
                                _serverErrors.remove('reportTypeDetails');
                              });
                            },
                            validator: (value) {
                              if (_isReporterCase && value == null) {
                                return 'Este campo es obligatorio';
                              }
                              return null;
                            },
                            serverError: _serverErrors['reportTypeDetails'],
                          ),
                          const SizedBox(height: 32),
                          _selectField(
                            label:
                                '¿Usted tiene algún rol de cuidado, representación o acompañamiento de la víctima? *',
                            placeholder:
                                'Indique si tiene rol de cuidado o representación',
                            value: _hasCareRole,
                            options: _yesNoOptions,
                            onChanged: (value) {
                              setState(() {
                                _hasCareRole = value;
                                _globalError = null;
                                _serverErrors.remove('hasCareRole');
                              });
                            },
                            validator: (value) {
                              if (_isReporterCase && value == null) {
                                return 'Este campo es obligatorio';
                              }
                              return null;
                            },
                            serverError: _serverErrors['hasCareRole'],
                          ),
                          const SizedBox(height: 32),
                          _selectField(
                            label:
                                '¿La víctima tiene conocimiento de este reporte?',
                            placeholder:
                                'Indique si la víctima conoce este reporte',
                            value: _victimAwareOfReport,
                            options: _yesNoOptions,
                            onChanged: (value) {
                              setState(() {
                                _victimAwareOfReport = value;
                                _globalError = null;
                                _serverErrors.remove('victimAwareOfReport');
                              });
                            },
                            validator: (_) => null,
                            serverError: _serverErrors['victimAwareOfReport'],
                          ),
                        ],
                        if (_shouldShowReporterIdentityFields) ...[
                          const SizedBox(height: 32),
                          _textField(
                            label: 'Nombres y apellidos de quién reporta',
                            controller: _reporterNamesCtrl,
                            hint:
                                'Por favor ingrese nombres y apellidos de quien reporta',
                            showEditIcon: true,
                            inputFormatters: _fullNameInputFormatters,
                            validator: (_) => null,
                            serverError: _serverErrors['reporterNames'],
                          ),
                          const SizedBox(height: 32),
                          _textField(
                            label: 'Teléfono de contacto de quién reporta',
                            controller: _reporterPhoneCtrl,
                            hint: '3001234567',
                            keyboardType: TextInputType.number,
                            showEditIcon: true,
                            inputFormatters: _colombiaPhoneInputFormatters,
                            prefixWidget: _colombiaPhonePrefix(),
                            validator: (value) => _validateColombiaPhone(
                              value,
                              required: false,
                            ),
                            serverError: _serverErrors['reporterPhone'],
                          ),
                        ],
                        if (_isVictimCase) ...[
                          _selectField(
                            label:
                                '¿Estaría dispuesto(a) a recibir una llamada para ahondar en la descripción de los hechos? *',
                            placeholder:
                                'Indique si estaría dispuesto(a) a recibir una llamada',
                            value: _willReceiveCall,
                            options: _yesNoOptions,
                            onChanged: (value) {
                              setState(() {
                                _willReceiveCall = value;
                                _globalError = null;
                                _serverErrors.remove('willReceiveCall');
                              });
                            },
                            validator: (value) {
                              if (_isVictimCase && value == null) {
                                return 'Este campo es obligatorio';
                              }
                              return null;
                            },
                            serverError: _serverErrors['willReceiveCall'],
                          ),
                        ],
                        const SizedBox(height: 32),
                        _textField(
                          label: 'Nombres de la víctima *',
                          controller: _victimNamesCtrl,
                          hint: 'Por favor ingrese nombres de la víctima',
                          showEditIcon: true,
                          inputFormatters: _nameInputFormatters,
                          validator: (value) {
                            if (!_showPersonalDataSection) return null;
                            if ((value ?? '').trim().isEmpty) {
                              return 'Este campo es obligatorio';
                            }
                            return null;
                          },
                          serverError: _serverErrors['names'],
                        ),
                        const SizedBox(height: 32),
                        _textField(
                          label: 'Apellidos de la víctima *',
                          controller: _victimLastNamesCtrl,
                          hint: 'Por favor ingrese apellidos de la víctima',
                          showEditIcon: true,
                          inputFormatters: _nameInputFormatters,
                          validator: (value) {
                            if (!_showPersonalDataSection) return null;
                            if ((value ?? '').trim().isEmpty) {
                              return 'Este campo es obligatorio';
                            }
                            return null;
                          },
                          serverError: _serverErrors['lastNames'],
                        ),
                        const SizedBox(height: 32),
                        _textField(
                          label:
                              'Télefono de contacto en Colombia de la víctima *',
                          controller: _victimColPhoneCtrl,
                          hint: '3001234567',
                          keyboardType: TextInputType.number,
                          showEditIcon: true,
                          inputFormatters: _colombiaPhoneInputFormatters,
                          prefixWidget: _colombiaPhonePrefix(),
                          validator: (value) {
                            if (!_showPersonalDataSection) return null;
                            return _validateColombiaPhone(
                              value,
                              required: true,
                            );
                          },
                          serverError: _serverErrors['victimColPhone'],
                        ),
                        const SizedBox(height: 32),
                        _textField(
                          label:
                              'Mejor hora de contacto para devolver la llamada a la víctima *',
                          controller: _bestContactTimeDisplayCtrl,
                          hint: 'Selecciona la hora (AM/PM)',
                          keyboardType: TextInputType.none,
                          showEditIcon: true,
                          readOnly: true,
                          onTap: _pickBestContactTime,
                          suffixWidget: const Icon(
                            Icons.schedule_rounded,
                            color: AppColors.salviaPurpura2,
                          ),
                          validator: (_) {
                            if (!_showPersonalDataSection) return null;
                            if (_bestContactTimeCtrl.text.trim().isEmpty) {
                              return 'Este campo es obligatorio';
                            }
                            return null;
                          },
                          serverError: _serverErrors['bestContactTime'],
                        ),
                        const SizedBox(height: 56),
                        _sectionHeader(
                          icon: Icons.description_outlined,
                          title: 'Descripción de los hechos',
                        ),
                        const SizedBox(height: 32),
                        _textField(
                          label: 'Descripción de los hechos',
                          controller: _factsCtrl,
                          hint: 'Por favor ingrese descripción de los hechos',
                          maxLines: 5,
                          showEditIcon: true,
                          contentPaddingOverride:
                              const EdgeInsets.fromLTRB(14, 28, 14, 14),
                          inputFormatters: const [
                            _WordLimitFormatter(maxWords: 500),
                          ],
                          maxWords: 500,
                          validator: (value) {
                            final words = _countWords(value ?? '');
                            if (words > 500) {
                              return 'Máximo 500 palabras';
                            }
                            return null;
                          },
                          serverError: _serverErrors['factsDescription'],
                        ),
                        if (!_skipCaptcha) ...[
                          const SizedBox(height: 32),
                          _buildCaptchaBlock(),
                        ],
                        const SizedBox(height: 24),
                        ElevatedButton(
                          onPressed: _isSubmitting ? null : _submit,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: AppColors.salviaTituloPrincipal,
                            disabledBackgroundColor: AppColors
                                .salviaTituloPrincipal
                                .withValues(alpha: 0.6),
                            foregroundColor: Colors.white,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(30),
                            ),
                            minimumSize: const Size(double.infinity, 82),
                          ),
                          child: _isSubmitting
                              ? const SizedBox(
                                  width: 18,
                                  height: 18,
                                  child: CircularProgressIndicator(
                                    strokeWidth: 2,
                                    color: Colors.white,
                                  ),
                                )
                              : const Row(
                                  mainAxisSize: MainAxisSize.min,
                                  children: [
                                    Text(
                                      'Enviar reporte',
                                      style: TextStyle(
                                        color: Colors.white,
                                        fontSize: 18,
                                        fontWeight: FontWeight.w400,
                                      ),
                                    ),
                                    SizedBox(width: 8),
                                    Icon(
                                      Icons.send_rounded,
                                      color: Colors.white,
                                      size: 24,
                                    ),
                                  ],
                                ),
                        ),
                      ],
                      espacio(30),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _sectionHeader({
    required IconData icon,
    required String title,
  }) {
    return Row(
      children: [
        Icon(
          icon,
          color: AppColors.salviaTituloPrincipal,
          size: 28,
        ),
        const SizedBox(width: 10),
        Expanded(
          child: Text(
            title,
            style: const TextStyle(
              color: AppColors.salviaTituloPrincipal,
              fontSize: _sectionHeaderFontSize,
              fontWeight: FontWeight.w700,
              letterSpacing: -0.7,
            ),
          ),
        ),
      ],
    );
  }

  Widget _paragraph(String text) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Text(
        text,
        textAlign: TextAlign.justify,
        style: const TextStyle(
          color: AppColors.salviaSubtitulos,
          fontSize: 14,
          height: 1.4,
        ),
      ),
    );
  }

  Widget _dataPolicyAccordion() {
    return Container(
      clipBehavior: Clip.antiAlias,
      decoration: BoxDecoration(
        color: const Color(0xfff7f5f0),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Theme(
        data: Theme.of(context).copyWith(dividerColor: Colors.transparent),
        child: ExpansionTile(
          key: const PageStorageKey<String>('data_policy_accordion'),
          backgroundColor: const Color(0xfff7f5f0),
          collapsedBackgroundColor: const Color(0xfff7f5f0),
          tilePadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 2),
          childrenPadding: const EdgeInsets.fromLTRB(14, 0, 14, 12),
          iconColor: AppColors.salviaPurpura2,
          collapsedIconColor: AppColors.salviaPurpura2,
          title: const Text(
            'Lectura y aceptación de la ley de protección de Datos',
            style: TextStyle(
              color: AppColors.salviaTitulos,
              fontSize: 16,
              fontWeight: FontWeight.w700,
            ),
          ),
          subtitle: const Text(
            'Toca para leer u ocultar el texto completo',
            style: TextStyle(
              color: AppColors.salviaPlaceholderColor,
              fontSize: 14,
              fontWeight: FontWeight.w400,
            ),
          ),
          children: [
            _paragraph(
              'Declaro de manera libre, expresa, inequívoca e informada, que AUTORIZO al Ministerio de Igualdad y Equidad para que de acuerdo con lo establecido en el literal a) del artículo 6 de la Ley 1581 de 2012, realice la recolección y tratamiento de mis datos personales que suministro de manera veraz y completa, los cuales serán utilizados para los diferentes aspectos asociados al Registro de Igualdad y Equidad, en cumplimiento de la Resolución 772 de 2024.',
            ),
            _paragraph(
              'A su vez, reconozco que se me ha informado de manera clara que tengo derecho a conocer, actualizar y rectificar los datos personales proporcionados, a solicitar prueba de esta autorización, a enterarme sobre el uso que se le ha dado a mis datos personales, a presentar quejas ante la Superintendencia de Industria y Comercio por el uso indebido de mis datos personales, a revocar esta autorización o solicitar la supresión de los datos personales suministrados y a acceder de forma gratuita a los mismos. Declaro que conozco y acepto la Política para el Tratamiento y Protección de Datos Personales del Ministerio de Igualdad y Equidad: Política para el Tratamiento y Protección de Datos Personales.',
            ),
            _paragraph(
              'Igualmente, manifiesto que de conformidad con el artículo 56 del Código de Procedimiento Administrativo y de lo Contencioso Administrativo - Ley 1437 de 2011 modificado por el artículo 10 de la Ley 2080 de 2021, autorizo expresamente al Ministerio de Igualdad y Equidad a remitir notificaciones electrónicas al correo electrónico proporcionado por mi parte.',
            ),
          ],
        ),
      ),
    );
  }

  Widget _errorBanner(String message) {
    return Container(
      margin: const EdgeInsets.only(bottom: 14),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: const Color(0xffffecec),
        border: Border.all(color: const Color(0xffdc3545)),
        borderRadius: BorderRadius.circular(10),
      ),
      child: Text(
        message,
        style: const TextStyle(color: Color(0xffdc3545), fontSize: 14),
      ),
    );
  }

  Widget _buildCaptchaBlock() {
    return Container(
      key: _captchaSectionKey,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Text(
                'Código de verificación',
                style: TextStyle(
                  color: AppColors.salviaTitulos,
                  fontSize: 16,
                  fontWeight: FontWeight.w700,
                ),
              ),
              const Spacer(),
              TextButton.icon(
                onPressed:
                    _isRefreshingCaptcha ? null : () => _refreshCaptchaId(),
                icon: _isRefreshingCaptcha
                    ? const SizedBox(
                        width: 14,
                        height: 14,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Icon(Icons.refresh, size: 18),
                label: const Text('Refrescar'),
              ),
            ],
          ),
          const SizedBox(height: 20),
          ClipRRect(
            borderRadius: BorderRadius.circular(12),
            child: _buildCaptchaPreview(),
          ),
          const SizedBox(height: 12),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: _isCaptchaAudioLoading ? null : _toggleCaptchaAudio,
              icon: _isCaptchaAudioLoading
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : Icon(
                      _isCaptchaAudioPlaying
                          ? Icons.stop_circle_outlined
                          : Icons.volume_up_rounded,
                    ),
              label: Text(
                _isCaptchaAudioPlaying
                    ? 'Detener audio'
                    : 'Escuchar código de verificación',
              ),
              style: OutlinedButton.styleFrom(
                foregroundColor: AppColors.salviaPurpura2,
                side: const BorderSide(color: AppColors.salviaPurpura4),
                padding: const EdgeInsets.symmetric(vertical: 24),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
            ),
          ),
          const SizedBox(height: 32),
          _textField(
            label: 'Escribe los números que ves en la imagen',
            controller: _captchaCtrl,
            hint: 'Escribe los números que ves en la imagen',
            keyboardType: TextInputType.number,
            inputFormatters: [
              FilteringTextInputFormatter.digitsOnly,
              LengthLimitingTextInputFormatter(6),
            ],
            showEditIcon: true,
            validator: (_) => null,
            serverError: _serverErrors['captchaSolution'],
          ),
        ],
      ),
    );
  }

  Widget _buildCaptchaPreview() {
    if (_isCaptchaLoading) {
      return Container(
        height: 68,
        width: double.infinity,
        alignment: Alignment.center,
        decoration: BoxDecoration(
          color: AppColors.salviaPurpura5,
          border: Border.all(color: AppColors.salviaPurpura4),
          borderRadius: BorderRadius.circular(12),
        ),
        child: const SizedBox(
          width: 18,
          height: 18,
          child: CircularProgressIndicator(strokeWidth: 2),
        ),
      );
    }

    if (_captchaBytes != null && _captchaBytes!.isNotEmpty) {
      return Container(
        height: 68,
        width: double.infinity,
        color: AppColors.salviaPurpura5,
        child: Image.memory(_captchaBytes!, fit: BoxFit.contain),
      );
    }

    return InkWell(
      onTap: _isRefreshingCaptcha ? null : () => _refreshCaptchaId(),
      child: Container(
        height: 68,
        width: double.infinity,
        alignment: Alignment.centerLeft,
        padding: const EdgeInsets.symmetric(horizontal: 26),
        decoration: BoxDecoration(
          color: AppColors.salviaPurpura5,
          border: Border.all(color: AppColors.salviaPurpura4),
          borderRadius: BorderRadius.circular(12),
        ),
        child: const Text(
          'El servidor de captcha devolvio imagen vacia. Toca para reintentar.',
          style: TextStyle(
            color: AppColors.salviaPlaceholderColor,
            fontSize: AppColors.salviaPlaceholderFontSize,
            fontWeight: FontWeight.w400,
          ),
          maxLines: 2,
          overflow: TextOverflow.ellipsis,
        ),
      ),
    );
  }

  Widget _selectField({
    required String label,
    required String placeholder,
    required _EnumOption? value,
    required List<_EnumOption> options,
    required ValueChanged<_EnumOption?> onChanged,
    required String? Function(_EnumOption?) validator,
    String? serverError,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            color: AppColors.salviaTitulos,
            fontSize: _fieldLabelFontSize,
            fontWeight: FontWeight.w700,
          ),
        ),
        const SizedBox(height: 8),
        DropdownButtonFormField<_EnumOption>(
          isExpanded: true,
          isDense: false,
          itemHeight: null,
          initialValue: value,
          hint: Text(
            placeholder,
            style: _formPlaceholderStyle,
          ),
          icon: const Icon(
            Icons.arrow_drop_down,
            color: AppColors.salviaPurpura2,
          ),
          decoration: InputDecoration(
            filled: true,
            fillColor: AppColors.salviaPurpura5,
            contentPadding:
                const EdgeInsets.symmetric(horizontal: 14, vertical: 16),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: AppColors.salviaPurpura4),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(
                color: AppColors.salviaPurpura2,
                width: 2,
              ),
            ),
            errorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Colors.red),
            ),
          ),
          validator: (selected) {
            final localError = validator(selected);
            if (localError != null) return localError;
            return serverError;
          },
          selectedItemBuilder: (BuildContext context) {
            return options.map((option) {
              return Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Container(
                    width: 8,
                    height: 8,
                    margin: const EdgeInsets.only(top: 7),
                    decoration: const BoxDecoration(
                      color: AppColors.salviaPurpura2,
                      shape: BoxShape.circle,
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Text(
                      option.name,
                      softWrap: true,
                      style: const TextStyle(
                        color: AppColors.salviaTitulos,
                        fontSize: 16,
                        fontWeight: FontWeight.w400,
                        height: 1.35,
                      ),
                    ),
                  ),
                ],
              );
            }).toList();
          },
          items: options
              .map(
                (option) => DropdownMenuItem<_EnumOption>(
                  value: option,
                  child: Padding(
                    padding: const EdgeInsets.symmetric(vertical: 12),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Container(
                          width: 8,
                          height: 8,
                          margin: const EdgeInsets.only(top: 7),
                          decoration: const BoxDecoration(
                            color: AppColors.salviaPurpura2,
                            shape: BoxShape.circle,
                          ),
                        ),
                        const SizedBox(width: 10),
                        Expanded(
                          child: Text(
                            option.name,
                            softWrap: true,
                            style: const TextStyle(
                              color: AppColors.salviaTitulos,
                              fontSize: 16,
                              fontWeight: FontWeight.w400,
                              height: 1.35,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              )
              .toList(),
          onChanged: onChanged,
        ),
      ],
    );
  }

  Widget _textField({
    required String label,
    required TextEditingController controller,
    required String hint,
    required String? Function(String?) validator,
    String? serverError,
    TextInputType keyboardType = TextInputType.text,
    int maxLines = 1,
    bool showEditIcon = false,
    List<TextInputFormatter>? inputFormatters,
    Widget? prefixWidget,
    Widget? suffixWidget,
    bool readOnly = false,
    VoidCallback? onTap,
    int? maxWords,
    EdgeInsetsGeometry? contentPaddingOverride,
  }) {
    final double verticalPadding = maxLines > 1 ? 14 : 18;
    final int currentWords =
        maxWords != null ? _countWords(controller.text) : 0;
    final bool highlightWordCounter = maxWords != null &&
        (maxWords == 500
            ? currentWords >= 450
            : currentWords >= (maxWords * 0.9).ceil());

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            if (showEditIcon) ...[
              const Icon(
                Icons.edit_outlined,
                size: 18,
                color: AppColors.salviaTituloPrincipal,
              ),
              const SizedBox(width: 8),
            ],
            Expanded(
              child: Text(
                label,
                style: const TextStyle(
                  color: AppColors.salviaTitulos,
                  fontSize: _fieldLabelFontSize,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        TextFormField(
          controller: controller,
          keyboardType: keyboardType,
          maxLines: maxLines,
          readOnly: readOnly,
          onTap: onTap,
          inputFormatters: inputFormatters,
          onChanged: (_) {
            final key = _serverErrorKeyFor(controller);
            setState(() {
              _globalError = null;
              if (key.isNotEmpty) {
                _serverErrors.remove(key);
              }
            });
          },
          decoration: InputDecoration(
            hintText: hint,
            hintStyle: _formPlaceholderStyle,
            prefixIcon: prefixWidget,
            suffixIcon: suffixWidget,
            prefixIconConstraints:
                const BoxConstraints(minWidth: 0, minHeight: 0),
            counterText:
                maxWords != null ? '$currentWords/$maxWords palabras' : null,
            counterStyle: TextStyle(
              color: highlightWordCounter
                  ? const Color(0xffdc3545)
                  : AppColors.salviaPlaceholderColor,
              fontSize: 12,
            ),
            filled: true,
            fillColor: AppColors.salviaPurpura5,
            contentPadding: contentPaddingOverride ??
                EdgeInsets.symmetric(
                  horizontal: 14,
                  vertical: verticalPadding,
                ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: AppColors.salviaPurpura4),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(
                color: AppColors.salviaPurpura2,
                width: 2,
              ),
            ),
            errorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Colors.red),
            ),
          ),
          validator: (value) {
            final localError = validator(value);
            if (localError != null) return localError;
            return serverError;
          },
        ),
      ],
    );
  }

  Widget _multiSelectField({
    required String label,
    required String placeholder,
    required List<_EnumOption> selectedOptions,
    required List<_EnumOption> options,
    required ValueChanged<List<_EnumOption>> onChanged,
    required String? Function(List<_EnumOption>) validator,
    String? serverError,
  }) {
    final valueKey = selectedOptions.map((e) => e.icode).join('|');

    return FormField<List<_EnumOption>>(
      key: ValueKey('multi_$valueKey'),
      initialValue: List<_EnumOption>.from(selectedOptions),
      validator: (value) {
        final active = value ?? const <_EnumOption>[];
        final localError = validator(active);
        if (localError != null) return localError;
        return serverError;
      },
      builder: (state) {
        final active = state.value ?? const <_EnumOption>[];
        final hasSelection = active.isNotEmpty;

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              label,
              style: const TextStyle(
                color: AppColors.salviaTitulos,
                fontSize: _fieldLabelFontSize,
                fontWeight: FontWeight.w700,
              ),
            ),
            const SizedBox(height: 8),
            InkWell(
              borderRadius: BorderRadius.circular(12),
              onTap: () async {
                final picked = await _showMultiSelectBottomSheet(
                  title: label.replaceAll('*', '').trim(),
                  options: options,
                  initialSelected: active,
                );
                if (picked == null) return;
                state.didChange(picked);
                onChanged(picked);
              },
              child: InputDecorator(
                isEmpty: !hasSelection,
                decoration: InputDecoration(
                  filled: true,
                  fillColor: AppColors.salviaPurpura5,
                  contentPadding:
                      const EdgeInsets.symmetric(horizontal: 14, vertical: 16),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide:
                        const BorderSide(color: AppColors.salviaPurpura4),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(
                      color: AppColors.salviaPurpura2,
                      width: 2,
                    ),
                  ),
                  errorBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(color: Colors.red),
                  ),
                  focusedErrorBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(color: Colors.red, width: 2),
                  ),
                  errorText: state.errorText,
                ),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Expanded(
                      child: hasSelection
                          ? Column(
                              crossAxisAlignment: CrossAxisAlignment.stretch,
                              children: [
                                for (var i = 0; i < active.length; i++)
                                  Padding(
                                    padding: EdgeInsets.only(
                                      bottom: i == active.length - 1 ? 0 : 8,
                                    ),
                                    child: _selectedAdjustmentChip(
                                      option: active[i],
                                      onRemove: () {
                                        final updated = List<_EnumOption>.from(
                                            active)
                                          ..removeWhere(
                                            (item) =>
                                                item.icode == active[i].icode,
                                          );
                                        state.didChange(updated);
                                        onChanged(updated);
                                      },
                                    ),
                                  ),
                              ],
                            )
                          : Text(
                              placeholder,
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                              style: _formPlaceholderStyle,
                            ),
                    ),
                    const SizedBox(width: 8),
                    Icon(
                      hasSelection
                          ? Icons.keyboard_arrow_up_rounded
                          : Icons.keyboard_arrow_down_rounded,
                      color: AppColors.salviaPurpura2,
                    ),
                  ],
                ),
              ),
            ),
          ],
        );
      },
    );
  }

  Widget _selectedAdjustmentChip({
    required _EnumOption option,
    required VoidCallback onRemove,
  }) {
    return Container(
      width: double.infinity,
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF6117D9), AppColors.salviaPurpura2],
          begin: Alignment.centerLeft,
          end: Alignment.centerRight,
        ),
        borderRadius: BorderRadius.circular(14),
      ),
      padding: const EdgeInsets.fromLTRB(12, 8, 8, 8),
      child: Row(
        mainAxisSize: MainAxisSize.max,
        children: [
          Expanded(
            child: Text(
              option.name,
              style: const TextStyle(
                color: Colors.white,
                fontSize: _fieldLabelFontSize,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          const SizedBox(width: 8),
          InkWell(
            borderRadius: BorderRadius.circular(12),
            onTap: onRemove,
            child: Container(
              width: 24,
              height: 24,
              decoration: BoxDecoration(
                color: Colors.white.withValues(alpha: 0.28),
                shape: BoxShape.circle,
              ),
              alignment: Alignment.center,
              child: const Icon(
                Icons.close_rounded,
                color: Colors.white,
                size: 16,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Future<List<_EnumOption>?> _showMultiSelectBottomSheet({
    required String title,
    required List<_EnumOption> options,
    required List<_EnumOption> initialSelected,
  }) async {
    final searchCtrl = TextEditingController();
    try {
      return await showModalBottomSheet<List<_EnumOption>>(
        context: context,
        isScrollControlled: true,
        useSafeArea: true,
        backgroundColor: Colors.transparent,
        builder: (context) {
          final tempSelected = List<_EnumOption>.from(initialSelected);
          String query = '';

          return StatefulBuilder(
            builder: (context, setSheetState) {
              final normalizedQuery = query.trim().toLowerCase();
              final filtered = options.where((opt) {
                if (normalizedQuery.isEmpty) return true;
                return opt.name.toLowerCase().contains(normalizedQuery);
              }).toList();

              return Container(
                margin: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: const Color(0xffffffff),
                  borderRadius: BorderRadius.circular(18),
                  border: Border.all(color: const Color(0xFF9E9E9E), width: 2),
                ),
                child: FractionallySizedBox(
                  heightFactor: 0.78,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      Padding(
                        padding: const EdgeInsets.fromLTRB(16, 14, 16, 8),
                        child: Text(
                          title,
                          style: const TextStyle(
                            color: AppColors.salviaTitulos,
                            fontSize: 16,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 16),
                        child: TextField(
                          controller: searchCtrl,
                          onChanged: (value) => setSheetState(() {
                            query = value;
                          }),
                          style: const TextStyle(
                            color: AppColors.salviaTitulos,
                            fontSize: 15,
                          ),
                          decoration: InputDecoration(
                            hintText: 'Buscar...',
                            hintStyle: _formPlaceholderStyle,
                            filled: true,
                            fillColor: AppColors.salviaPurpura5,
                            contentPadding: const EdgeInsets.symmetric(
                              horizontal: 14,
                              vertical: 14,
                            ),
                            enabledBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: const BorderSide(
                                color: AppColors.salviaPurpura2,
                                width: 1.6,
                              ),
                            ),
                            focusedBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: const BorderSide(
                                color: AppColors.salviaPurpura2,
                                width: 2,
                              ),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 8),
                      const Divider(height: 1, color: Color(0xFF9E9E9E)),
                      Expanded(
                        child: ListView.builder(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 8,
                            vertical: 6,
                          ),
                          itemCount: filtered.length,
                          itemBuilder: (context, index) {
                            final option = filtered[index];
                            final checked = _containsOption(
                              tempSelected,
                              option,
                            );

                            return CheckboxListTile(
                              controlAffinity: ListTileControlAffinity.leading,
                              value: checked,
                              activeColor: AppColors.salviaPurpura2,
                              checkColor: Colors.white,
                              side: const BorderSide(
                                color: Color(0xFF9E9E9E),
                                width: 1.2,
                              ),
                              contentPadding: const EdgeInsets.symmetric(
                                horizontal: 8,
                                vertical: 2,
                              ),
                              title: Text(
                                option.name,
                                style: const TextStyle(
                                  color: AppColors.salviaTitulos,
                                  fontSize: 15,
                                  fontWeight: FontWeight.w400,
                                  height: 1.35,
                                ),
                              ),
                              onChanged: (_) {
                                setSheetState(() {
                                  _toggleOption(tempSelected, option);
                                });
                              },
                            );
                          },
                        ),
                      ),
                      const Divider(height: 1, color: Color(0xFF9E9E9E)),
                      Padding(
                        padding: const EdgeInsets.fromLTRB(10, 8, 10, 12),
                        child: Row(
                          children: [
                            TextButton(
                              onPressed: () {
                                setSheetState(() {
                                  tempSelected.clear();
                                });
                              },
                              child: const Text('Limpiar'),
                            ),
                            const Spacer(),
                            TextButton(
                              onPressed: () {
                                Navigator.of(context).pop();
                              },
                              child: const Text('Cancelar'),
                            ),
                            const SizedBox(width: 8),
                            ElevatedButton(
                              onPressed: () {
                                Navigator.of(
                                  context,
                                ).pop(List<_EnumOption>.from(tempSelected));
                              },
                              style: ElevatedButton.styleFrom(
                                backgroundColor: AppColors.salviaPurpura2,
                                foregroundColor: Colors.white,
                              ),
                              child: const Text('Aplicar'),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          );
        },
      );
    } finally {
      searchCtrl.dispose();
    }
  }

  bool _containsOption(List<_EnumOption> selected, _EnumOption option) {
    for (final item in selected) {
      if (item.icode == option.icode) return true;
    }
    return false;
  }

  void _toggleOption(List<_EnumOption> selected, _EnumOption option) {
    final index = selected.indexWhere((item) => item.icode == option.icode);
    if (index >= 0) {
      selected.removeAt(index);
      return;
    }
    selected.add(option);
  }

  String _serverErrorKeyFor(TextEditingController controller) {
    if (controller == _reporterNamesCtrl) return 'reporterNames';
    if (controller == _reporterPhoneCtrl) return 'reporterPhone';
    if (controller == _victimNamesCtrl) return 'names';
    if (controller == _victimLastNamesCtrl) return 'lastNames';
    if (controller == _victimColPhoneCtrl) return 'victimColPhone';
    if (controller == _bestContactTimeCtrl ||
        controller == _bestContactTimeDisplayCtrl) {
      return 'bestContactTime';
    }
    if (controller == _factsCtrl) return 'factsDescription';
    if (controller == _captchaCtrl) return 'captchaSolution';
    return '';
  }

  String _onlyDigits(String value) => value.replaceAll(RegExp(r'\D'), '');
  int? _digitsToInt(String value) {
    final digits = _onlyDigits(value);
    if (digits.isEmpty) return null;
    return int.tryParse(digits);
  }

  int _countWords(String value) => _countWordsInText(value);

  Future<void> _pickBestContactTime() async {
    final initial =
        _parse24HourTime(_bestContactTimeCtrl.text.trim()) ?? TimeOfDay.now();

    final picked = await showTimePicker(
      context: context,
      initialTime: initial,
      builder: (context, child) {
        return MediaQuery(
          data: MediaQuery.of(context).copyWith(alwaysUse24HourFormat: false),
          child: child ?? const SizedBox.shrink(),
        );
      },
    );

    if (picked == null || !mounted) return;

    setState(() {
      _bestContactTimeCtrl.text = _formatTime24(picked);
      _bestContactTimeDisplayCtrl.text = _formatTime12(picked);
      _globalError = null;
      _serverErrors.remove('bestContactTime');
    });
  }

  String _formatTime24(TimeOfDay time) {
    final hh = time.hour.toString().padLeft(2, '0');
    final mm = time.minute.toString().padLeft(2, '0');
    return '$hh:$mm';
  }

  String _formatTime12(TimeOfDay time) {
    final hour12 = time.hourOfPeriod == 0 ? 12 : time.hourOfPeriod;
    final mm = time.minute.toString().padLeft(2, '0');
    final suffix = time.period == DayPeriod.am ? 'AM' : 'PM';
    return '$hour12:$mm $suffix';
  }

  TimeOfDay? _parse24HourTime(String value) {
    final match = RegExp(r'^([01]?\d|2[0-3]):([0-5]\d)$').firstMatch(value);
    if (match == null) return null;
    final hour = int.tryParse(match.group(1) ?? '');
    final minute = int.tryParse(match.group(2) ?? '');
    if (hour == null || minute == null) return null;
    return TimeOfDay(hour: hour, minute: minute);
  }

  Widget _colombiaPhonePrefix() {
    return const Padding(
      padding: EdgeInsets.only(left: 14, right: 8),
      child: Text(
        '🇨🇴 +57',
        style: TextStyle(
          color: AppColors.salviaTitulos,
          fontSize: 15,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  String? _validateColombiaPhone(
    String? value, {
    required bool required,
  }) {
    final digits = _onlyDigits(value ?? '');
    if (digits.isEmpty) {
      return required ? 'Este campo es obligatorio' : null;
    }
    if (digits.length != 10) {
      return 'Ingresa un número móvil de 10 dígitos';
    }
    return null;
  }

  Future<void> _loadInitialFormAndCaptcha() async {
    final config = await ServerProxy().fetchPublicFormConfig();
    if (!mounted) return;

    if (config != null) {
      final remoteYesNo = _toEnumOptions(config.yesNo);
      final remoteReportType = _toEnumOptions(config.reportType);
      final remoteReportTypeDetails = _toEnumOptions(config.reportTypeDetails);
      final remoteAdjustmentsGBV = _toEnumOptions(config.adjustmentsGBV);
      debugPrint(
        '[form-config] yesNo=${remoteYesNo.length} reportType=${remoteReportType.length} reportTypeDetails=${remoteReportTypeDetails.length} adjustmentsGBV=${remoteAdjustmentsGBV.length} captchaID="${config.captchaId}"',
      );

      setState(() {
        if (remoteYesNo.isNotEmpty) {
          _yesNoOptions = remoteYesNo;
        }
        if (remoteReportType.isNotEmpty) {
          _reportTypeOptions = remoteReportType;
        }
        if (remoteReportTypeDetails.isNotEmpty) {
          _reportTypeDetailsOptions = remoteReportTypeDetails;
        }
        if (remoteAdjustmentsGBV.isNotEmpty) {
          _adjustmentsGBVOptions = remoteAdjustmentsGBV;
        }

        _authorizationAnswer =
            _remapOption(_authorizationAnswer, _yesNoOptions);
        _hasCareRole = _remapOption(_hasCareRole, _yesNoOptions);
        _victimAwareOfReport =
            _remapOption(_victimAwareOfReport, _yesNoOptions);
        _willReceiveCall = _remapOption(_willReceiveCall, _yesNoOptions);
        _reportType = _remapOption(_reportType, _reportTypeOptions);
        _reportTypeDetails =
            _remapOption(_reportTypeDetails, _reportTypeDetailsOptions);
        _adjustmentsGBVSelected = _remapSelectedOptions(
          _adjustmentsGBVSelected,
          _adjustmentsGBVOptions,
        );

        if (!_skipCaptcha && config.captchaId.trim().isNotEmpty) {
          _captchaId = config.captchaId.trim();
          _captchaCacheBust = DateTime.now().millisecondsSinceEpoch;
          _captchaCtrl.clear();
        }
      });
    }

    if (_skipCaptcha) return;

    if (config != null && config.captchaId.trim().isNotEmpty) {
      await _loadCaptchaImage(showErrorMessage: false);
      return;
    }

    await _refreshCaptchaId(showErrorMessage: false);
  }

  Future<void> _refreshCaptchaId({bool showErrorMessage = true}) async {
    await _captchaAudioPlayer.stop();
    setState(() {
      _isRefreshingCaptcha = true;
      _isCaptchaAudioPlaying = false;
    });

    final fetchedId = await ServerProxy().fetchCaptchaIdFromPublicForm();
    if (!mounted) return;

    if (fetchedId != null && fetchedId.trim().isNotEmpty) {
      setState(() {
        _captchaId = fetchedId.trim();
        _captchaCacheBust = DateTime.now().millisecondsSinceEpoch;
        _captchaCtrl.clear();
        _isRefreshingCaptcha = false;
      });
      await _loadCaptchaImage(showErrorMessage: showErrorMessage);
      return;
    }

    setState(() {
      _captchaCacheBust = DateTime.now().millisecondsSinceEpoch;
      _isRefreshingCaptcha = false;
    });

    if (showErrorMessage) {
      final reason = fetchedId == null ? 'respuesta nula' : 'captchaID vacío';
      const webHint = kIsWeb ? ' (en web puede ser CORS del servidor)' : '';
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            'No se pudo obtener captcha ($reason)$webHint. Base: ${ServerProxy.baseUrl}',
          ),
          backgroundColor: Colors.red,
        ),
      );
    }

    await _loadCaptchaImage(showErrorMessage: false);
  }

  Future<void> _loadCaptchaImage({bool showErrorMessage = true}) async {
    setState(() {
      _isCaptchaLoading = true;
    });

    final bytes = await ServerProxy().fetchCaptchaImageBytes(
      _captchaId,
      cacheBust: _captchaCacheBust,
    );
    if (!mounted) return;

    setState(() {
      _captchaBytes = bytes;
      _isCaptchaLoading = false;
    });

    if ((bytes == null || bytes.isEmpty) && showErrorMessage) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text(
            'El servidor devolvio el captcha vacio. Intenta refrescar en unos segundos.',
          ),
          backgroundColor: Colors.orange,
        ),
      );
    }
  }

  Future<void> _toggleCaptchaAudio() async {
    if (_isCaptchaAudioLoading) return;
    if (_isCaptchaAudioPlaying) {
      await _captchaAudioPlayer.stop();
      return;
    }

    setState(() {
      _isCaptchaAudioLoading = true;
    });

    try {
      final cacheBust = DateTime.now().millisecondsSinceEpoch;
      final url = ServerProxy.captchaAudioUrl(_captchaId, cacheBust: cacheBust);
      debugPrint('[captcha] AUDIO $url');

      final audioBytes = await ServerProxy().fetchCaptchaAudioBytes(
        _captchaId,
        cacheBust: cacheBust,
      );

      await _captchaAudioPlayer.stop();
      await _captchaAudioPlayer.setVolume(1.0);
      if (audioBytes != null && audioBytes.isNotEmpty) {
        await _captchaAudioPlayer.play(
          BytesSource(audioBytes, mimeType: 'audio/wav'),
        );
      } else {
        // Fallback por URL en caso de que por alguna razón fallen los bytes.
        await _captchaAudioPlayer.play(UrlSource(url));
      }
    } on MissingPluginException {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text(
            'Audio captcha no disponible en esta sesión. Reinicia la app completa.',
          ),
          backgroundColor: Colors.orange,
        ),
      );
    } catch (_) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text(
            kIsWeb
                ? 'No se pudo reproducir el audio captcha en este navegador.'
                : 'No se pudo reproducir el audio captcha.',
          ),
          backgroundColor: Colors.orange,
        ),
      );
    } finally {
      if (mounted) {
        setState(() {
          _isCaptchaAudioLoading = false;
        });
      }
    }
  }

  Future<void> _resolveCurrentLocation() async {
    try {
      final serviceEnabled = await Geolocator.isLocationServiceEnabled();
      if (!serviceEnabled) {
        debugPrint(
          '[location] location services disabled, using fallback lat/lng',
        );
        return;
      }

      var permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        permission = await Geolocator.requestPermission();
      }

      if (permission == LocationPermission.denied ||
          permission == LocationPermission.deniedForever) {
        debugPrint(
          '[location] location permission denied, using fallback lat/lng',
        );
        return;
      }

      const locationSettings = LocationSettings(
        accuracy: LocationAccuracy.medium,
        timeLimit: Duration(seconds: 8),
      );
      final position = await Geolocator.getCurrentPosition(
        locationSettings: locationSettings,
      );

      if (!mounted) return;

      setState(() {
        _latitude = position.latitude;
        _longitude = position.longitude;
      });

      debugPrint(
        '[location] using device lat=$_latitude lng=$_longitude (web=$kIsWeb)',
      );
    } catch (e) {
      debugPrint(
        '[location] could not obtain device location, using fallback: $e',
      );
    }
  }

  List<_EnumOption> _toEnumOptions(List<Map<String, dynamic>> rawOptions) {
    final out = <_EnumOption>[];
    for (final item in rawOptions) {
      final icode = (item['icode'] ?? '').toString().trim();
      final name = (item['name'] ?? '').toString().trim();
      final code = (item['code'] ?? '').toString().trim();
      if (icode.isEmpty || name.isEmpty || code.isEmpty) {
        continue;
      }
      out.add(_EnumOption(icode: icode, name: name, code: code));
    }
    return out;
  }

  _EnumOption? _remapOption(_EnumOption? current, List<_EnumOption> options) {
    if (current == null) return null;
    for (final opt in options) {
      if (opt.icode == current.icode) return opt;
    }
    for (final opt in options) {
      if (opt.code == current.code) return opt;
    }
    return null;
  }

  List<_EnumOption> _remapSelectedOptions(
    List<_EnumOption> selected,
    List<_EnumOption> options,
  ) {
    if (selected.isEmpty) return <_EnumOption>[];
    final remapped = <_EnumOption>[];
    for (final item in selected) {
      final mapped = _remapOption(item, options);
      if (mapped != null && !_containsOption(remapped, mapped)) {
        remapped.add(mapped);
      }
    }
    return remapped;
  }

  Future<void> _submit() async {
    final isValid = _formKey.currentState?.validate() ?? false;
    if (!isValid) {
      _scrollController.animateTo(
        0,
        duration: const Duration(milliseconds: 350),
        curve: Curves.easeOut,
      );
      return;
    }

    if (!_showPersonalDataSection) return;

    setState(() {
      _isSubmitting = true;
      _globalError = null;
      _serverErrors.clear();
    });

    final reporterPhone = _digitsToInt(_reporterPhoneCtrl.text);
    final victimColPhone = _digitsToInt(_victimColPhoneCtrl.text);
    if (victimColPhone == null) {
      setState(() {
        _isSubmitting = false;
        _globalError = 'El número de teléfono de la víctima no es válido.';
      });
      return;
    }

    final payload = {
      'victimContact': {
        'form2': {
          'authorizationAnswer': _authorizationAnswer?.toJson(),
          if (_isAuthorized)
            'adjustmentsGBV': _adjustmentsGBVSelected
                .map((option) => option.toJson())
                .toList(),
          'reportType': _reportType?.toJson(),
          if (_isReporterCase)
            'reportTypeDetails': _reportTypeDetails?.toJson(),
          if (_isReporterCase) 'hasCareRole': _hasCareRole?.toJson(),
          if (_isReporterCase && _victimAwareOfReport != null)
            'victimAwareOfReport': _victimAwareOfReport?.toJson(),
          if (_shouldShowReporterIdentityFields)
            'reporterNames': _reporterNamesCtrl.text.trim(),
          if (_shouldShowReporterIdentityFields && reporterPhone != null)
            'reporterPhone': reporterPhone,
          if (_isVictimCase) 'willReceiveCall': _willReceiveCall?.toJson(),
          'victimColPhone': victimColPhone,
          'bestContactTime': _bestContactTimeCtrl.text.trim(),
          if (_factsCtrl.text.trim().isNotEmpty)
            'factsDescription': _factsCtrl.text.trim(),
        },
        'names': _victimNamesCtrl.text.trim(),
        'lastNames': _victimLastNamesCtrl.text.trim(),
        'latitude': _latitude,
        'longitude': _longitude,
        'captchaID': _skipCaptcha ? '' : _captchaId,
        'captchaSolution': _skipCaptcha ? '' : _captchaCtrl.text.trim(),
      }
    };
    debugPrint('[submit] payload=${jsonEncode(payload)}');

    final sp = ServerProxy();
    final res = await sp.submitPrimerContacto(payload);
    debugPrint(
      '[submit] status=${res.statusCode} success=${res.success} body=${res.mensaje}',
    );

    if (!mounted) return;

    _captchaCtrl.clear();

    if (res.success) {
      if (!mounted) return;
      Navigator.pushAndRemoveUntil(
        context,
        MaterialPageRoute(builder: (_) => const ProcesandoReporteScreen()),
        (Route<dynamic> route) => false,
      );
    } else {
      _applyServerErrors(res);
    }

    setState(() {
      _isSubmitting = false;
    });
  }

  void _applyServerErrors(RespuestaServidor res) {
    final Map<String, dynamic> data = res.data;
    final Map<String, String> collected = {};

    final defaultErrors = _asMap(data['default']);
    final topVictim = _asMap(data['victimContact']);
    final topForm2 = _asMap(data['form2']);
    final dottedForm2 = _asMap(data['victimContact.form2']);
    final nestedForm2 = _asMap(topVictim['form2']);

    final form2Errors = topForm2.isNotEmpty
        ? topForm2
        : (nestedForm2.isNotEmpty ? nestedForm2 : dottedForm2);

    if (defaultErrors['global_msg'] != null) {
      _globalError = _normalizeUserMessage(
        defaultErrors['global_msg'].toString(),
      );
    } else if (res.mensaje.isNotEmpty) {
      _globalError = _normalizeUserMessage(res.mensaje);
    }

    if (defaultErrors['captchaID'] != null) {
      _captchaId = defaultErrors['captchaID'].toString();
      _captchaCacheBust = DateTime.now().millisecondsSinceEpoch;
      _loadCaptchaImage(showErrorMessage: false);
    }

    _captureError(form2Errors, 'authorizationAnswer', collected);
    _captureError(form2Errors, 'adjustmentsGBV', collected);
    _captureError(form2Errors, 'reportType', collected);
    _captureError(form2Errors, 'reportTypeDetails', collected);
    _captureError(form2Errors, 'hasCareRole', collected);
    _captureError(form2Errors, 'victimAwareOfReport', collected);
    _captureError(form2Errors, 'reporterNames', collected);
    _captureError(form2Errors, 'reporterPhone', collected);
    _captureError(form2Errors, 'willReceiveCall', collected);
    _captureError(form2Errors, 'victimColPhone', collected);
    _captureError(form2Errors, 'bestContactTime', collected);
    _captureError(form2Errors, 'factsDescription', collected);

    _captureError(topVictim, 'names', collected);
    _captureError(topVictim, 'lastNames', collected);
    _captureError(topVictim, 'captchaSolution', collected);
    if (collected['captchaSolution'] != null) {
      collected['captchaSolution'] = _normalizeUserMessage(
        collected['captchaSolution']!,
      );
    }
    final captchaErrorMessage = collected['captchaSolution'];

    final topVictimDefault = (topVictim['default'] ?? '').toString().trim();
    if (topVictimDefault.isNotEmpty &&
        !_isGenericServerMessage(topVictimDefault)) {
      _globalError = topVictimDefault;
    }

    if (collected.isEmpty &&
        (_globalError == null ||
            _globalError!.trim().isEmpty ||
            _isGenericServerMessage(_globalError!))) {
      final fallbackMessages = _extractErrorMessages(data);
      if (fallbackMessages.isNotEmpty) {
        _globalError = fallbackMessages.take(4).join('\n');
      }
    }

    if ((_globalError == null || _isGenericServerMessage(_globalError!)) &&
        res.mensaje.trim().isNotEmpty &&
        !_isGenericServerMessage(res.mensaje)) {
      _globalError = res.mensaje;
    }

    if (captchaErrorMessage != null && captchaErrorMessage.trim().isNotEmpty) {
      _globalError = captchaErrorMessage;
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!mounted) return;
        final ctx = _captchaSectionKey.currentContext;
        if (ctx != null) {
          Scrollable.ensureVisible(
            ctx,
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOut,
            alignment: 0.1,
          );
        }
        ScaffoldMessenger.of(context).hideCurrentSnackBar();
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            behavior: SnackBarBehavior.floating,
            backgroundColor: Colors.orange,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(14),
            ),
            content: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Icon(
                  Icons.warning_amber_rounded,
                  color: Colors.white,
                  size: 22,
                ),
                const SizedBox(width: 10),
                Expanded(
                  child: Text(
                    captchaErrorMessage,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                      fontWeight: FontWeight.w400,
                      height: 1.3,
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      });
    }

    setState(() {
      _serverErrors
        ..clear()
        ..addAll(collected);
    });

    _formKey.currentState?.validate();
  }

  Map<String, dynamic> _asMap(dynamic raw) {
    if (raw is Map<String, dynamic>) return raw;
    if (raw is Map) return Map<String, dynamic>.from(raw);
    return {};
  }

  void _captureError(
    Map<String, dynamic> source,
    String key,
    Map<String, String> output,
  ) {
    final value = source[key];
    if (value == null) return;

    if (value is List && value.isNotEmpty) {
      output[key] = value.first.toString();
      return;
    }

    output[key] = value.toString();
  }

  List<String> _extractErrorMessages(dynamic source) {
    final messages = <String>[];

    void visit(String path, dynamic value) {
      if (value == null) return;

      if (value is String) {
        final text = value.trim();
        if (text.isEmpty) return;
        if (_isGenericServerMessage(text)) {
          return;
        }
        messages.add(path.isEmpty ? text : '$path: $text');
        return;
      }

      if (value is List) {
        for (final item in value) {
          visit(path, item);
        }
        return;
      }

      if (value is Map) {
        for (final entry in value.entries) {
          final key = entry.key.toString();
          if (key == 'global_msg' || key == 'captchaID') continue;
          final nextPath = path.isEmpty ? key : '$path.$key';
          visit(nextPath, entry.value);
        }
      }
    }

    visit('', source);

    final seen = <String>{};
    final unique = <String>[];
    for (final msg in messages) {
      if (seen.add(msg)) unique.add(msg);
    }
    return unique;
  }

  bool _isGenericServerMessage(String text) {
    final lower = text.trim().toLowerCase();
    return lower == 'se presentaron los siguientes errores' ||
        lower == 'se produjo un error general' ||
        lower == 'se ha producido un error inesperado';
  }

  String _normalizeUserMessage(String text) {
    return text.replaceAll(RegExp(r'captcha', caseSensitive: false), 'CÓDIGO');
  }
}

class _EnumOption {
  final String icode;
  final String name;
  final String code;

  const _EnumOption({
    required this.icode,
    required this.name,
    required this.code,
  });

  Map<String, String> toJson() => {
        'icode': icode,
        'name': name,
        'code': code,
      };
}

class _NameTextFormatter extends TextInputFormatter {
  final int maxLength;

  const _NameTextFormatter({required this.maxLength});

  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    final raw = newValue.text;
    if (raw.isEmpty) return newValue;

    // Mientras el teclado está componiendo (subrayado), no alteramos el texto
    // para no romper tildes ni backspace.
    if (!newValue.composing.isCollapsed) return newValue;

    var sanitized = raw.replaceAll(RegExp(r'[^A-Za-zÀ-ÖØ-öø-ÿ\s]'), '');
    if (sanitized.length > maxLength) {
      sanitized = sanitized.substring(0, maxLength);
    }

    final lowered = sanitized.toLowerCase();
    final StringBuffer out = StringBuffer();
    bool capitalizeNext = true;
    final letterRegex = RegExp(r'[a-zà-öø-ÿ]');

    for (var i = 0; i < lowered.length; i++) {
      final ch = lowered[i];
      if (capitalizeNext && letterRegex.hasMatch(ch)) {
        out.write(ch.toUpperCase());
        capitalizeNext = false;
      } else {
        out.write(ch);
        if (ch.trim().isNotEmpty && letterRegex.hasMatch(ch)) {
          capitalizeNext = false;
        }
      }

      if (ch == ' ') {
        capitalizeNext = true;
      }
    }

    final formatted = out.toString();
    if (formatted == newValue.text) return newValue;

    final rawOffset = newValue.selection.baseOffset;
    final safeOffset = rawOffset < 0
        ? formatted.length
        : (rawOffset > formatted.length ? formatted.length : rawOffset);

    return TextEditingValue(
      text: formatted,
      selection: TextSelection.collapsed(offset: safeOffset),
      composing: newValue.composing,
    );
  }
}

int _countWordsInText(String text) {
  final trimmed = text.trim();
  if (trimmed.isEmpty) return 0;
  return trimmed.split(RegExp(r'\s+')).where((w) => w.isNotEmpty).length;
}

class _WordLimitFormatter extends TextInputFormatter {
  final int maxWords;

  const _WordLimitFormatter({required this.maxWords});

  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    if (!newValue.composing.isCollapsed) return newValue;

    final oldCount = _countWordsInText(oldValue.text);
    final newCount = _countWordsInText(newValue.text);
    if (newCount <= maxWords) return newValue;

    // Permite retroceder y editar sin bloquear cuando ya se alcanzó el tope.
    if (newValue.text.length < oldValue.text.length) return newValue;
    if (newCount <= oldCount) return newValue;

    return oldValue;
  }
}
