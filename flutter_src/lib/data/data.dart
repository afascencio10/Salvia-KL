// data.dart
class Option {
  String title;
  String description;
  int pendingMessages;

  Option(
      {required this.title,
      required this.description,
      required this.pendingMessages});
}

final List<Option> options = [
  Option(
      title:
          'Corta descripción del caso que reportó el usuario cuando usó la aplicación.',
      description: 'Leer más',
      pendingMessages: 2),
  Option(
      title:
          'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque pulvinar porta nibh at tempor. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque pulvinar porta nibh at tempor.',
      description: 'Leer más',
      pendingMessages: 5),
];
