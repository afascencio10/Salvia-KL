//import 'package:demo_flut/global/utils.dart';
import 'package:flutter/material.dart';
import '../../global/my_colors.dart';

class DropReportar extends StatefulWidget {
  final double fontSize;
  final Color fontColor;
  final IconData dropdownIcon;
  final double iconSize;
  final Color iconColor;
  final List<String> opciones;
  final Function(String) onSeleccion;

  const DropReportar({
    Key? key,
    this.fontSize = 16.0,
    this.fontColor = const Color(0xff45008b),
    this.dropdownIcon = Icons.arrow_drop_down,
    this.iconSize = 30.0,
    this.iconColor = const Color(0xff8C48BF),
    required this.opciones,
    required this.onSeleccion,
  }) : super(key: key);

  @override
  DropReportarState createState() => DropReportarState();
}

class DropReportarState extends State<DropReportar> {
  String? dropdownValue;

  @override
  void initState() {
    super.initState();
    if (widget.opciones.isNotEmpty) {
      dropdownValue = widget.opciones.first;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.start,
      children: <Widget>[
        Expanded(
          child: Container(
            padding: const EdgeInsets.symmetric(vertical: 2, horizontal: 16),
            decoration: BoxDecoration(
              color: AppColors.salviaPurpura5, // Color de fondo rosado claro
              border: Border.all(color: AppColors.salviaPurpura4, width: 2),
              borderRadius: BorderRadius.circular(20),
            ),
            child: DropdownButtonHideUnderline(
              child: DropdownButton<String>(
                isExpanded: true,
                value: dropdownValue,
                onChanged: (String? newValue) {
                  setState(() {
                    dropdownValue = newValue;
                    widget.onSeleccion(newValue!);
                  });
                },
                items: widget.opciones
                    .map<DropdownMenuItem<String>>((String value) {
                  return DropdownMenuItem<String>(
                    value: value,
                    child: Text(
                      value,
                      style: TextStyle(
                        fontSize: widget.fontSize,
                        color: widget.fontColor,
                      ),
                    ),
                  );
                }).toList(),
                icon: Icon(
                  widget.dropdownIcon,
                  color: widget.iconColor,
                  size: widget.iconSize,
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }
}
