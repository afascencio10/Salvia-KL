import 'dart:async';
import 'package:flutter/material.dart';
import 'package:record/record.dart';
import 'package:audioplayers/audioplayers.dart';
import '../../global/my_colors.dart';

class AudioRecorderPlayer extends StatefulWidget {
  const AudioRecorderPlayer({Key? key}) : super(key: key);

  @override
  AudioRecorderPlayerState createState() => AudioRecorderPlayerState();
}

class AudioRecorderPlayerState extends State<AudioRecorderPlayer> {
  final Record _record = Record();
  final AudioPlayer _audioPlayer = AudioPlayer();
  bool _isRecording = false;
  bool _isPlaying = false;
  String? _filePath;
  Duration _duration = Duration.zero;
  Duration _recordingDuration = Duration.zero;
  Timer? _timer;
  Timer? _recordingTimer;

  @override
  void initState() {
    super.initState();
    _audioPlayer.onPlayerComplete.listen((event) {
      _stopPlayback();
    });
  }

  void _startTimer() {
    _timer = Timer.periodic(const Duration(seconds: 1), (Timer t) async {
      Duration? positionDuration = await _audioPlayer.getCurrentPosition();
      int position = positionDuration?.inMilliseconds ?? 0;
      setState(() {
        _duration = Duration(milliseconds: position);
      });
    });
  }

  void _stopTimer() {
    _timer?.cancel();
    setState(() {
      _duration = Duration.zero;
    });
  }

  void _startRecordingTimer() {
    _recordingTimer = Timer.periodic(const Duration(seconds: 1), (timer) {
      setState(() {
        _recordingDuration += const Duration(seconds: 1);
      });
    });
  }

  void _stopRecordingTimer() {
    _recordingTimer?.cancel();
  }

  void _toggleRecording() async {
    if (_isRecording) {
      _filePath = await _record.stop();
      _stopRecordingTimer();
      setState(() {
        _isRecording = false;
      });
    } else {
      // Aquí, antes de comenzar la grabación, restablece el contador.
      setState(() {
        _recordingDuration = Duration.zero;
      });
      await _record.start();
      _startRecordingTimer();
      setState(() {
        _isRecording = true;
      });
    }
  }

  void _togglePlayback() async {
    if (_isPlaying) {
      _stopPlayback();
    } else {
      if (_filePath != null) {
        await _audioPlayer.play(UrlSource(_filePath!));
        _startTimer();
        setState(() {
          _isPlaying = true;
        });
      }
    }
  }

  void _stopPlayback() {
    _audioPlayer.stop();
    _stopTimer();
    setState(() {
      _isPlaying = false;
      _duration = Duration.zero; // Resetear la duración de reproducción
    });
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.start,
      children: <Widget>[
        const SizedBox(width: 15),
        // Botón de grabación
        Container(
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: _isRecording
                ? AppColors.salviaPurpura2
                : AppColors.salviaPurpura1,
            borderRadius: BorderRadius.circular(12),
          ),
          child: IconButton(
            icon: Icon(
              _isRecording ? Icons.stop : Icons.mic,
              color: Colors.white,
            ),
            onPressed: _toggleRecording,
          ),
        ),
        const SizedBox(width: 8),
        // Botón de reproducción
        Container(
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: _isPlaying
                ? AppColors.salviaTurquesa2
                : AppColors.salviaTurquesa1,
            borderRadius: BorderRadius.circular(12),
          ),
          child: IconButton(
            icon: Icon(
              _isPlaying ? Icons.stop : Icons.play_arrow,
              color: Colors.white,
            ),
            onPressed: _filePath == null ? null : _togglePlayback,
          ),
        ),
        const SizedBox(width: 12),
        // Contadores en columna
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Contador de grabación
            Text(
              "Rec: ${_recordingDuration.inMinutes}:${_recordingDuration.inSeconds.remainder(60).toString().padLeft(2, '0')}",
              style: const TextStyle(
                color: AppColors.salviaPurpura1,
                fontSize: 16.0,
              ),
            ),
            const SizedBox(height: 12),
            // Contador de reproducción
            Text(
              "Rep: ${_duration.inMinutes}:${_duration.inSeconds.remainder(60).toString().padLeft(2, '0')}",
              style: const TextStyle(
                color: AppColors.salviaTurquesa1,
                fontSize: 16.0,
              ),
            ),
          ],
        ),
      ],
    );
  }
}
