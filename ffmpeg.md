# Build ffmpeg
```
https://ffmpeg.org/releases/ffmpeg-snapshot.tar.bz2
tar xf ffmpeg-snapshot.tar.bz2
cd ffmpeg/
./configure --enable-debug=3 --disable-optimizations --disable-asm
```

# Debug
```
gdb ffmpeg_g
GNU gdb (GDB) Fedora Linux 13.2-6.fc38
Copyright (C) 2023 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
Type "show copying" and "show warranty" for details.
This GDB was configured as "x86_64-redhat-linux-gnu".
Type "show configuration" for configuration details.
For bug reporting instructions, please see:
<https://www.gnu.org/software/gdb/bugs/>.
Find the GDB manual and other documentation resources online at:
    <http://www.gnu.org/software/gdb/documentation/>.

For help, type "help".
Type "apropos word" to search for commands related to "word"...
Reading symbols from ffmpeg_g...
(gdb) list
1284	    proc = GetCurrentProcess();
1285	    memcounters.cb = sizeof(memcounters);
1286	    GetProcessMemoryInfo(proc, &memcounters, sizeof(memcounters));
1287	    return memcounters.PeakPagefileUsage;
1288	#else
1289	    return 0;
1290	#endif
1291	}
1292	
1293	int main(int argc, char **argv)
(gdb) break main
Breakpoint 1 at 0x43f944: file fftools/ffmpeg.c, line 1298.
(gdb) run
Starting program: /home/knightingal/source/ffmpeg/ffmpeg_g 

This GDB supports auto-downloading debuginfo from the following URLs:
  <https://debuginfod.fedoraproject.org/>
Enable debuginfod for this session? (y or [n]) y
Debuginfod has been enabled.
To make this setting permanent, add 'set debuginfod enabled on' to .gdbinit.
Downloading separate debug info for system-supplied DSO at 0x7ffff7fc7000
Downloading separate debug info for /lib64/libm.so.6                                                                                                                                          
Downloading separate debug info for /lib64/libxcb.so.1                                                                                                                                        
Downloading separate debug info for /lib64/libxcb-shm.so.0                                                                                                                                    
Downloading separate debug info for /home/knightingal/.cache/debuginfod_client/959bbf9bd1fee9880f5599adb988633bb6794612/debuginfo                                                             
Downloading separate debug info for /lib64/libxcb-shape.so.0                                                                                                                                  
Downloading separate debug info for /lib64/libxcb-xfixes.so.0                                                                                                                                 
Downloading separate debug info for /lib64/libbz2.so.1                                                                                                                                        
Downloading separate debug info for /home/knightingal/.cache/debuginfod_client/f8237ff24b622047428b97b1486b393e3a9e38de/debuginfo                                                             
Downloading separate debug info for /lib64/libz.so.1                                                                                                                                          
Downloading separate debug info for /lib64/liblzma.so.5                                                                                                                                       
Downloading separate debug info for /lib64/libc.so.6                                                                                                                                          
[Thread debugging using libthread_db enabled]                                                                                                                                                 
Using host libthread_db library "/lib64/libthread_db.so.1".
Downloading separate debug info for /lib64/libXau.so.6
Breakpoint 1, main (argc=1, argv=0x7fffffffde28) at fftools/ffmpeg.c:1298
1298	    init_dynload();
Missing separate debuginfos, use: dnf debuginfo-install glibc-2.37-13.fc38.x86_64 libXau-1.0.11-2.fc38.x86_64
(gdb) n
1300	    setvbuf(stderr,NULL,_IONBF,0); /* win32 runtime needs this */
(gdb) 
1302	    av_log_set_flags(AV_LOG_SKIP_REPEATED);
(gdb) 
1303	    parse_loglevel(argc, argv, options);
(gdb) 
1306	    avdevice_register_all();
(gdb) 
1308	    avformat_network_init();
(gdb) 
1310	    show_banner(argc, argv, options);
(gdb) 
ffmpeg version N-112863-gea6817d2a7 Copyright (c) 2000-2023 the FFmpeg developers
  built with gcc 13 (GCC)
  configuration: --enable-debug=3 --disable-optimizations --disable-asm
  libavutil      58. 32.100 / 58. 32.100
  libavcodec     60. 35.100 / 60. 35.100
  libavformat    60. 18.100 / 60. 18.100
  libavdevice    60.  4.100 / 60.  4.100
  libavfilter     9. 13.100 /  9. 13.100
  libswscale      7.  6.100 /  7.  6.100
  libswresample   4. 13.100 /  4. 13.100
1313	    ret = ffmpeg_parse_options(argc, argv);
(gdb) 
1314	    if (ret < 0)
(gdb) 
1317	    if (nb_output_files <= 0 && nb_input_files == 0) {
(gdb) print ret
$1 = 0
(gdb) 
$2 = 0
(gdb) 
$3 = 0
(gdb) n
1318	        show_usage();
(gdb) 
Hyper fast Audio and Video encoder
usage: ffmpeg [options] [[infile options] -i infile]... {[outfile options] outfile}...

1319	        av_log(NULL, AV_LOG_WARNING, "Use -h to get full help or, even better, run 'man %s'\n", program_name);
(gdb) 
Use -h to get full help or, even better, run 'man ffmpeg'
1320            ret = 1;
(gdb) 
1321	        goto finish;
(gdb) 
1347	    if (ret == AVERROR_EXIT)
(gdb) print ret
$4 = 1
(gdb) n
1350	    ffmpeg_cleanup(ret);
(gdb) 
1351	    return ret;
(gdb) 
1352	}
(gdb) 
0x00007ffff7c72b8a in __libc_start_call_main () from /lib64/libc.so.6
(gdb) 
Single stepping until exit from function __libc_start_call_main,
which has no line number information.
[Inferior 1 (process 16172) exited with code 01]
(gdb) exit
```