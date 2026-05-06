import 'package:flutter/material.dart';
import '../../global/my_colors.dart';

class MyBadge extends StatelessWidget {
  final int messages;

  const MyBadge({Key? key, required this.messages}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(6),
      decoration: BoxDecoration(
        color: AppColors.salviaAmarillo1,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        messages.toString(),
        style: const TextStyle(color: Colors.white, fontSize: 16),
      ),
    );
  }
}
