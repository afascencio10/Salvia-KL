import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../../global/my_colors.dart';

enum PickerMode { date, time, both }

class IconoYDateTimePicker extends StatefulWidget {
  final double iconSize;
  final Color iconColor;
  final String hintText;
  final PickerMode pickerMode;
  final Color dateTextColor;
  final Color timeTextColor;
  final Function(DateTime)? onDateSelected;
  final Function(TimeOfDay)? onTimeSelected;

  const IconoYDateTimePicker({
    Key? key,
    this.iconSize = 50.0,
    this.iconColor = AppColors.salviaPurpura1,
    this.hintText = 'Seleccionar fecha y hora',
    required this.pickerMode,
    this.dateTextColor = AppColors.salviaPurpura1,
    this.timeTextColor = AppColors.salviaPurpura1,
    this.onDateSelected,
    this.onTimeSelected,
  }) : super(key: key);

  @override
  IconoYDateTimePickerState createState() => IconoYDateTimePickerState();
}

class IconoYDateTimePickerState extends State<IconoYDateTimePicker> {
  DateTime? selectedDate;
  TimeOfDay? selectedTime;

  Future<void> _pickDateTime() async {
    // ------------------------------------------
    // Lógica para la FECHA
    // ------------------------------------------
    ThemeData datePickerTheme = ThemeData.light().copyWith(
      colorScheme: ColorScheme.light(
        primary: widget.dateTextColor,
        onPrimary: Colors.white,
      ),
      dialogTheme: const DialogThemeData(backgroundColor: Colors.white),
    );

    if (widget.pickerMode == PickerMode.date ||
        widget.pickerMode == PickerMode.both) {
      final DateTime? date = await showDatePicker(
        context: context,
        initialDate: selectedDate ?? DateTime.now(),
        firstDate: DateTime(2025),
        lastDate: DateTime(2035),
        builder: (BuildContext context, Widget? child) {
          return Theme(
            data: datePickerTheme,
            child: child!,
          );
        },
      );

      // Verificamos si el widget sigue montado después de la espera
      if (!mounted) return;

      if (date != null) {
        selectedDate = date;
        widget.onDateSelected?.call(date);
        setState(() {});
      }
    }

    // ------------------------------------------
    // Lógica para la HORA
    // ------------------------------------------
    ThemeData timePickerTheme = ThemeData.light().copyWith(
      colorScheme: ColorScheme.light(
        primary: widget.timeTextColor,
        onPrimary: Colors.white,
      ),
      dialogTheme: const DialogThemeData(backgroundColor: Colors.white),
    );

    if (widget.pickerMode == PickerMode.time ||
        widget.pickerMode == PickerMode.both) {
      // CORRECCIÓN 2: Verificamos mounted antes de usar 'context' en showTimePicker
      // Esto soluciona el "async gap" si venimos de elegir una fecha.
      if (!mounted) return;

      final TimeOfDay? time = await showTimePicker(
        context: context,
        initialTime: selectedTime ?? TimeOfDay.now(),
        builder: (BuildContext context, Widget? child) {
          return Theme(
            data: timePickerTheme,
            child: child!,
          );
        },
      );

      if (!mounted) return;

      if (time != null) {
        selectedTime = time;
        widget.onTimeSelected?.call(time);
        setState(() {});
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    TextSpan dateText = selectedDate != null
        ? TextSpan(
            text: DateFormat.yMd().format(selectedDate!),
            style: TextStyle(color: widget.dateTextColor),
          )
        : const TextSpan(text: '');

    TextSpan timeText = selectedTime != null
        ? TextSpan(
            text: selectedTime!.format(context),
            style: TextStyle(color: widget.timeTextColor),
          )
        : const TextSpan(text: '');

    IconData iconData;
    switch (widget.pickerMode) {
      case PickerMode.date:
        iconData = Icons.calendar_today;
        break;
      case PickerMode.time:
        iconData = Icons.access_time;
        break;
      case PickerMode.both: // CORRECCIÓN 1: Eliminamos 'default'
        iconData = Icons.date_range;
        break;
    }

    return Row(
      mainAxisAlignment: MainAxisAlignment.start,
      children: <Widget>[
        IconButton(
          icon: Icon(iconData),
          onPressed: _pickDateTime,
          color: widget.iconColor,
          iconSize: widget.iconSize,
        ),
        const SizedBox(width: 20.0),
        Expanded(
          child: RichText(
            text: TextSpan(
              style: const TextStyle(color: Colors.black),
              children: [dateText, timeText],
            ),
          ),
        ),
      ],
    );
  }
}
