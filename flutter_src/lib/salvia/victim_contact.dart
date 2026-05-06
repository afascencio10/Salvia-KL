/// Modelo VictimContact – todos los campos son String.
class VictimContact {
  // -------------------- Atributos --------------------
  final String statusDescription;
  final String names;
  final String lastNames;
  final String nick;
  final String docType;
  final String docNumber;
  final String birthDate; // ISO‑8601 (yyyy-MM-ddTHH:mm:ssZ)
  final String townCode;
  final String address;
  final String livingLatitude; // se guarda como string
  final String livingLongitude; // se guarda como string
  final String phone;
  final String genderIdentity;
  final String sexualOrientation;
  final String origin;
  final String occupation;
  final String occupationOther;
  final String factsDescription;

  // --- NUEVOS CAMPOS ---
  final String contactDate; // Fecha para contactar
  final String contactTime; // Hora para contactar

  // -------------------- Constructor --------------------
  VictimContact({
    this.statusDescription = '',
    this.names = '',
    this.lastNames = '',
    this.nick = '',
    this.docType = '',
    this.docNumber = '',
    this.birthDate = '',
    this.townCode = '',
    this.address = '',
    this.livingLatitude = '',
    this.livingLongitude = '',
    this.phone = '',
    this.genderIdentity = '',
    this.sexualOrientation = '',
    this.origin = '',
    this.occupation = '',
    this.occupationOther = '',
    this.factsDescription = '',

    // --- Inicializar nuevos campos ---
    this.contactDate = '',
    this.contactTime = '',
  });

  // -------------------- Deserialización --------------------
  factory VictimContact.fromJson(Map<String, dynamic> json) {
    return VictimContact(
      statusDescription: _toString(json['statusDescription']),
      names: _toString(json['names']),
      lastNames: _toString(json['lastNames']),
      nick: _toString(json['nick']),
      docType: _toString(json['docType']),
      docNumber: _toString(json['docNumber']),
      birthDate: _dateToIso8601(json['birthDate']),
      townCode: _toString(json['townCode']),
      address: _toString(json['address']),
      livingLatitude: _numToString(json['livingLatitude']),
      livingLongitude: _numToString(json['livingLongitude']),
      phone: _toString(json['phone']),
      genderIdentity: _toString(json['genderIdentity']),
      sexualOrientation: _toString(json['sexualOrientation']),
      origin: _toString(json['origin']),
      occupation: _toString(json['occupation']),
      occupationOther: _toString(json['occupationOther']),
      factsDescription: _toString(json['factsDescription']),
      // --- Leer nuevos campos ---
      contactDate: _toString(json['contactDate']),
      contactTime: _toString(json['contactTime']),
    );
  }

  // -------------------- Serialización --------------------
  Map<String, dynamic> toJson() {
    return {
      if (statusDescription.isNotEmpty) 'statusDescription': statusDescription,
      if (names.isNotEmpty) 'names': names,
      if (lastNames.isNotEmpty) 'lastNames': lastNames,
      if (nick.isNotEmpty) 'nick': nick,
      if (docType.isNotEmpty) 'docType': docType,
      if (docNumber.isNotEmpty) 'docNumber': docNumber,
      // Se mantiene el formato ISO‑8601
      if (birthDate.isNotEmpty)
        'birthDate': birthDate.isNotEmpty ? _parseIso8601(birthDate) : '',
      if (townCode.isNotEmpty) 'townCode': townCode,
      if (address.isNotEmpty) 'address': address,
      if (livingLatitude.isNotEmpty) 'livingLatitude': livingLatitude,
      if (livingLongitude.isNotEmpty) 'livingLongitude': livingLongitude,
      if (phone.isNotEmpty) 'phone': phone,
      if (genderIdentity.isNotEmpty) 'genderIdentity': genderIdentity,
      if (sexualOrientation.isNotEmpty) 'sexualOrientation': sexualOrientation,
      if (origin.isNotEmpty) 'origin': origin,
      if (occupation.isNotEmpty) 'occupation': occupation,
      if (occupationOther.isNotEmpty) 'occupationOther': occupationOther,
      if (factsDescription.isNotEmpty) 'factsDescription': factsDescription,
      // --- Enviar nuevos campos ---
      if (contactDate.isNotEmpty) 'contactDate': contactDate,
      if (contactTime.isNotEmpty) 'contactTime': contactTime,
    };
  }

  // -------------------- Helpers de conversión --------------------
  /// Convierte cualquier valor a String, devolviendo '' si es null.
  static String _toString(dynamic value) => value?.toString() ?? '';

  /// Convierte una fecha (DateTime, String ISO‑8601 o nulo) a string ISO‑8601.
  static String _dateToIso8601(dynamic value) {
    if (value == null || value.toString().trim().isEmpty) return '';
    try {
      final dt = value is DateTime ? value : DateTime.parse(value.toString());
      // Normalizamos a ISO‑8601 sin milisegundos si no los tiene.
      return dt.toUtc().toIso8601String();
    } catch (_) {
      // Si la cadena no se puede parsear, devolvemos ''.
      return '';
    }
  }

  /// Convierte una fecha ISO‑8601 string a DateTime y vuelve a string
  /// (para que quede en el mismo formato al serializar).
  static String _parseIso8601(String value) {
    try {
      final dt = DateTime.parse(value);
      return dt.toUtc().toIso8601String();
    } catch (_) {
      return '';
    }
  }

  /// Convierte un número (int, double o string con dígitos) a string.
  static String _numToString(dynamic value) => value?.toString() ?? '';
}
