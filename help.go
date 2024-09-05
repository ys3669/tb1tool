// help.go
package main

import (
	"fmt"

	"github.com/fatih/color"
)

func Showh() {
	fmt.Println("Available Commands and Options: 書きかけ")
	fmt.Println("-p    : Serial port to use (default: /dev/ttyUSB0)")
	fmt.Println("-s    : Baud rate for the serial port (default: 38400)")
	fmt.Println("-g    : API,CFG,SYS to generate NMEA sentences")
	fmt.Println("-z    : API,CFG,SYS with parameters to send via serial port")
	fmt.Println("-S    : Send NMEA sentence to update baud rate")
	fmt.Println("-h	   : Display help message")
	fmt.Println("--help: Display long help message")
}

// ShowHelp displays the available configuration options and detailed command information
func ShowHelp() {
	yellow := color.New(color.FgYellow).Add(color.Underline).SprintFunc()
	cyan := color.New(color.FgCyan).Add(color.Underline).SprintFunc()
	Showh()

	fmt.Println("\nDetailed Command Information:")

	fmt.Println("\n7.1 API [GNSS] – 受信する衛星の種類を選択する")
	fmt.Println("  - GET: GNSS,QUERY")
	fmt.Println("  - POST: GNSS,TalkerID,Gps,Glonass,Galileo,Qzss,Sbas/L1s")

	fmt.Println("\n7.2 API [PPS] – PPS に関する各種設定を行う")
	fmt.Println("  - POST: PPS,Type,Mode,Period,Pulsewidth,Cable delay,Polarity")

	fmt.Println("\n7.3 API [SURVEY] – 位置モードの設定を行う")
	fmt.Println("  - POST: SURVEY,Position mode[Sigma threshold,Time threshold][Latitude,Longitude,Altitude]")

	fmt.Println("\n7.4 API [RESTART] – リスタートを実施する")
	fmt.Println("  - POST: RESTART[Restart type]")

	fmt.Println("\n7.5 API [FLASHBACKUP] – FLASH ROM へのバックアップを実施する")
	fmt.Println("  - GET: FLASHBACKUP,QUERY")
	fmt.Println("  - POST: FLASHBACKUP,Type")

	fmt.Println("\n7.6 API [DEFLS] – デフォルトうるう秒を設定する")
	fmt.Println("  - GET: DEFLS,QUERY")
	fmt.Println("  - POST: DEFLS,Sec")

	fmt.Println("\n7.7 API [TIMEZONE] – LZT を設定する")
	fmt.Println("  - POST: TIMEZONE,Sign,Hour,Minute[sec]")

	fmt.Println("\n7.8 API [TIMEALIGN] – 出力時刻と PPS の同期対象を設定する")
	fmt.Println("  - GET: TIMEALIGN,QUERY")
	fmt.Println("  - POST: TIMEALIGN,Mode")

	fmt.Println("\n7.9 API [TIME] – 時刻を設定する")
	fmt.Println("  - POST: TIME,Time of day,Day,Month,Year")

	fmt.Println("\n7.10 API [FIXMASK] – 受信衛星に対して各種マスクを設定する")
	fmt.Println("  - GET: FIXMASK,QUERY")
	fmt.Println("  - POST: FIXMASK,Mode,Elevmask,Reserve,SNRmask,IDSM[Prohibit SVs(GPS),Prohibit SVs(GLONASS),Prohibit SVs(Galileo),Prohibit SVs(QZSS),Prohibit SVs(SBAS)]")

	fmt.Println("\n7.11 API [OCP] – 詳細な仰角・方位角マスクを設定する")
	fmt.Println("  - POST: OCP,num,el_0,el_1,…,el_19")
	fmt.Println("  - POST: OCP,az_1,el_1[az_2,el_2[ … [Az_9,El_9]…]")

	fmt.Println("\n7.12 API [NLOSMASK] – NLOS 衛星を排除するアルゴリズムの各種設定をする")
	fmt.Println("  - GET: NLOSMASK,QUERY")
	fmt.Println("  - POST: NLOSMASK,mode,Threshold1,Threshold2,Threshold3")

	fmt.Println("\n7.13 API [MODESET] – 周波数モードの状態遷移に関する閾値を設定する")
	fmt.Println("  - GET: MODESET,QUERY")
	fmt.Println("  - POST: MODESET,Lock port,Coarse lock threshold,phase skip threshold")

	fmt.Println("\n7.14 API [PHASESKIP] – 位相飛ばしフラグを設定する")
	fmt.Println("  - POST: PHASESKIP,phase skip flag")

	fmt.Println("\n7.15 API [HOSET] – ホールドオーバーの学習時間と実施時間に関する設定する")
	fmt.Println("  - GET: HOSET,QUERY")
	fmt.Println("  - POST: HOSET,Ho set flag[Learning time set0,Available time set0,Learning time set1,Available time set1,Learning time set2,Available time set2]")

	fmt.Println("\n7.16 API [ANTSET] – アンテナ給電に関する設定する")
	fmt.Println("  - GET: ANTSET,QUERY")
	fmt.Println("  - POST: ANTSET,Antenna status")

	fmt.Println("\n7.17 API [ALMSET] – アラーム出力に関する設定する")
	fmt.Println("  - GET: ALMSET,QUERY")
	fmt.Println("  - POST: ALMSET,Force alarm,Alarm mask")

	fmt.Println("\n7.18 API [CROUT] – CR センテンスの出力状態を設定する")
	fmt.Println("  - POST: CROUT,Type,Rate")

	fmt.Println("===========================================================================")

	fmt.Println("7.19", yellow("CFG"), "[NMEAOUT] – 標準 NMEA の出力状態を設定する")
	fmt.Println("  - POST: NMEAOUT,Type,Interval")

	fmt.Println("7.20", yellow("CFG"), "[UART1] – シリアルポートの設定")
	fmt.Println("  - POST: UART1,Baudrate")

	fmt.Println("7.21", cyan("SYS"), "[VERSION] –ソフトウェアバージョンの出力要求")
	fmt.Println("  - GET: VERSION")
}
