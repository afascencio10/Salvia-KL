
import 'package:flutter/material.dart';
import '../../global/my_colors.dart';

class AudioButton extends StatelessWidget {
  final int buttonNumber;
  final String audioName;
  final bool isPlaying;
  final int? playingButton;
  final Function(String, int) onToggleAudio;

  const AudioButton({super.key, 
    required this.buttonNumber,
    required this.audioName,
    required this.isPlaying,
    required this.playingButton,
    required this.onToggleAudio,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: playingButton == buttonNumber && isPlaying
            ? AppColors.salviaPurpura4
            : AppColors.salviaPurpura5,
        borderRadius: BorderRadius.circular(50),
      ),
      child: IconButton(
        icon: Icon(
          playingButton == buttonNumber && isPlaying
              ? Icons.stop
              : Icons.play_arrow,
          color: playingButton == buttonNumber && isPlaying
              ? AppColors.salviaPurpura2
              : AppColors.salviaPurpura2,
        ),
        iconSize: 30.0,
        onPressed: () => onToggleAudio(audioName, buttonNumber),
      ),
    );
  }
}
