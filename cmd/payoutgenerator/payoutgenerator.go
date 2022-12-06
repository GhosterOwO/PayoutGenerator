package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"
)

type User struct {
	Common struct {
		Gameid  string   `json:"gameid"`
		Symbols []string `json:"Symbols"`
	} `json:"common"`
	Mg struct {
		Field       []int   `json:"field"`
		Stripe90    [][]int `json:"stripe90"`
		Paylines    [][]int `json:"paylines"`
		Paytable    [][]int `json:"paytable"`
		Quapaytable [][]int `json:"quapaytable"`
	} `json:"mg"`
	Fg1 struct {
		Field       []int   `json:"field"`
		Stripe90    [][]int `json:"stripe90"`
		Paylines    [][]int `json:"paylines"`
		Paytable    [][]int `json:"paytable"`
		Quapaytable [][]int `json:"quapaytable"`
	} `json:"fg1"`
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

// 定義變數----------------------------------------------
var spintimes int = 0             //Spin次數
var costmoney int = 0             //累積押注金額
var mgtotalwin int = 0            //Mg總贏分
var mgsctotalwin int = 0          //MgSC總贏分
var fgtotalwin int = 0            //Fg總贏分
var fgtotalwin3 int = 0           //Fg總贏分
var fgtotalwin4 int = 0           //Fg總贏分
var fgtotalwin5 int = 0           //Fg總贏分
var mgrtp float64 = 0             //MgRTP
var fgrtp float64 = 0             //FgRTP
var fgrtp3 float64 = 0            //FgRTP
var fgrtp4 float64 = 0            //FgRTP
var fgrtp5 float64 = 0            //FgRTP
var mgScrtp float64 = 0           //FgRTP
var roundwin int = 0              //回合總贏分
var fghit3 float64 = 0            //進Fg次數
var fghit4 float64 = 0            //進Fg次數
var fghit5 float64 = 0            //進Fg次數
var fghit3percent float64 = 0     //進Fg頻率百分比
var fghit4percent float64 = 0     //進Fg頻率百分比
var fghit5percent float64 = 0     //進Fg頻率百分比
var fgretrigerhit3 float64 = 0    //進Fg次數
var fgretrigerhit4 float64 = 0    //進Fg次數
var fgretrigerhit5 float64 = 0    //進Fg次數
var fgretrigerhit3per float64 = 0 //進Fg頻率百分比
var fgretrigerhit4per float64 = 0 //進Fg頻率百分比
var fgretrigerhit5per float64 = 0 //進Fg頻率百分比
var rtp float64 = 0               //遊戲總體RTP
var hitpercent float64 = 0        //遊戲hitrate百分比
var winpercent float64 = 0        //遊戲winrate百分比
var cv float64                    //遊戲波動率
var kkk int = 0                   //FreeSpin總轉數
var user User                     //讀取資料
var SymbolWin [29][5]int          //各Symbol分數LOG
var SymbolWin2 [29][5]float64     //各Symbol分數LOG

// Round 四捨五入，ROUND_HALF_UP 模式實現
// 返回將 val 根據指定精度 precision（十進位制小數點後數字的數目）進行四捨五入的結果。precision 也可以是負數或零
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}

// 畫面停止按任意鍵離開
func pause() {
	fmt.Printf("按任意键退出...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

// 顯示進度條
func loadber(i int) {
	if i%10000000 == 0 && i != 0 {
		fmt.Print("\n*")
	} else if i%1000000 == 0 {
		fmt.Print("*")
	}
}

// 創建一個錯誤處理函數，避免過多的 if err != nil{} 出現
func dropErr(e error) {
	if e != nil {
		fmt.Println("開啟失敗")
		panic(e)
	}
}

func transpose(A [][]string) {
	// 交換行和列索引
	result := make([][]string, len(A[0]))
	for i, _ := range result {
		result[i] = make([]string, len(A))
	}
	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A[0]); j++ {
			result[j][i] = A[i][j]
		}
	}
	fmt.Println("\n")
	for j := 0; j < len(A[0]); j++ {
		fmt.Println(result[j])
	}
}

func max(a []int) int {
	var i, k int
	k = a[0]
	for i = 1; i < len(a); i++ {
		if a[i] > k {
			k = a[i]
		}
	}
	return k
}

func loadfile() {
	// 打開文件
	jsonFile, err := os.Open("./set01.json")
	dropErr(err)
	// 關閉文件
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &user)
}

func Outputfile() {
	fmt.Println("\n計算結果如下：\n花費:", costmoney, "元", "\t\trtp:", Round(rtp*100, 3), "%", "\nmg獲得:", mgtotalwin, "元", "\t\tfg獲得:", fgtotalwin, "元", "\nmg_RTP:", Round(mgrtp*100, 3), "%", "\t\tfg_RTP:", Round(fgrtp*100, 3), "%", "\nhitrate:", Round(hitpercent*100, 3), "%", "\t\twinrate:", Round(winpercent*100, 3), "%")
	fmt.Println("\n進FG次數：", fghit3+fghit4+fghit5, "\t\t進FG機率:", Round((fghit3percent+fghit4percent+fghit5percent)*100, 3), "%", "\nFGRetriger次數", fgretrigerhit3+fgretrigerhit4+fgretrigerhit5, "\t\tFGRetriger機率:", Round((fgretrigerhit3per+fgretrigerhit4per+fgretrigerhit5per)*100, 3), "\ncv:", Round(cv, 3), "%")
	fmt.Println("\nFGhit%(3):", Round(fghit3percent*100, 3), "%", "\t\tFGretrigerhit%(3):", Round(fgretrigerhit3per*100, 3), "%", "\nFGhit%(4):", Round(fghit4percent*100, 3), "%", "\t\tFGretrigerhit%(4):", Round(fgretrigerhit4per*100, 3), "%", "\nFGhit%(5):", Round(fghit5percent*100, 3), "%", "\t\tFGretrigerhit%(5):", Round(fgretrigerhit5per*100, 3), "%")
	fmt.Println("\nFG3個:", fgtotalwin3, "元", "\t\tfg_3個RTP:", Round(fgrtp3*100, 3), "%", "\nFG4個:", fgtotalwin4, "元", "\t\tfg_4個RTP:", Round(fgrtp4*100, 3), "%", "\nFG5個:", fgtotalwin5, "元", "\t\tfg_5個RTP:", Round(fgrtp5*100, 3), "%")
	//fmt.Println("\n\n H1:", SymbolWin2[indexOf("H1", user.Common.Symbols)], "\n H2:", SymbolWin2[indexOf("H2", user.Common.Symbols)], "\n H3:", SymbolWin2[indexOf("H3", user.Common.Symbols)], "\n L1:", SymbolWin2[indexOf("L1", user.Common.Symbols)], "\n L2:", SymbolWin2[indexOf("L2", user.Common.Symbols)], "\n L3:", SymbolWin2[indexOf("L3", user.Common.Symbols)], "\n L4:", SymbolWin2[indexOf("L4", user.Common.Symbols)], "\n L5:", SymbolWin2[indexOf("L5", user.Common.Symbols)], "\n W1:", SymbolWin2[indexOf("W1", user.Common.Symbols)], "\n F1:", SymbolWin2[indexOf("F1", user.Common.Symbols)])

	fmt.Println("\n成功輸出output.txt檔案至資料夾中")
	fmt.Print("\n")

	// 打開文件將資訊寫入
	file, err := os.OpenFile("output.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0604)
	dropErr(err)
	fmt.Fprint(file, "計算結果如下：\n花費:", costmoney, "元", "\trtp:", Round(rtp*100, 4), "%", "\nmg獲得:", mgtotalwin, "元", "\tfg獲得:", fgtotalwin, "元", "\nmg_RTP:", Round(mgrtp*100, 4), "%", "\tfg_RTP:", Round(fgrtp*100, 4), "%", "\nhitrate:", Round(hitpercent*100, 6), "%", "\twinrate:", Round(winpercent*100, 6), "%\n") // 向file对应文件中写入数据
	fmt.Fprint(file, "\n進FG次數", fghit3+fghit4+fghit5, "\t進FG機率:", Round((fghit3percent+fghit4percent+fghit5percent)*100, 4), "%", "\nFGRetriger次數", fgretrigerhit3+fgretrigerhit4+fgretrigerhit5, "\tFGRetriger機率:", Round((fgretrigerhit3per+fgretrigerhit4per+fgretrigerhit5per)*100, 4), "\ncv:", Round(cv, 6), "%\n")
	fmt.Fprint(file, "\nFGhit%(3):", Round(fghit3percent*100, 4), "%", "\tFGretrigerhit%(3):", Round(fgretrigerhit3per*100, 4), "%", "\nFGhit%(4):", Round(fghit4percent*100, 4), "%", "\tFGretrigerhit%(4):", Round(fgretrigerhit4per*100, 4), "%", "\nFGhit%(5):", Round(fghit5percent*100, 4), "%", "\tFGretrigerhit%(5):", Round(fgretrigerhit5per*100, 4), "%")
	fmt.Fprint(file, "\n\n H1:", SymbolWin2[indexOf("H1", user.Common.Symbols)], "\n H2:", SymbolWin2[indexOf("H2", user.Common.Symbols)], "\n H3:", SymbolWin2[indexOf("H3", user.Common.Symbols)], "\n L1:", SymbolWin2[indexOf("L1", user.Common.Symbols)], "\n L2:", SymbolWin2[indexOf("L2", user.Common.Symbols)], "\n L3:", SymbolWin2[indexOf("L3", user.Common.Symbols)], "\n L4:", SymbolWin2[indexOf("L4", user.Common.Symbols)], "\n L5:", SymbolWin2[indexOf("L5", user.Common.Symbols)], "\n W1:", SymbolWin2[indexOf("W1", user.Common.Symbols)], "\n F1:", SymbolWin2[indexOf("F1", user.Common.Symbols)])
	file.Close()
	pause()
}

// 開始執行主程式-------------------------------------------
func main() {
	loadfile()
	fmt.Println("請輸入旋轉次數:")
	fmt.Scanln(&spintimes)
	mgSpin(spintimes)
	Outputfile()
}

// MainGame流程--------------------------------------------
func mgSpin(times int) {
	//MG轉盤流程
	rand.Seed(time.Now().Unix())
	var i, j, h int
	var hitp, winp float64 //hit & win計算
	var pos [][]int        //盤面
	var ptmp []int         //盤面參數
	for i = 0; i < len(user.Mg.Field); i++ {
		for j = 0; j < user.Mg.Field[i]; j++ {
			ptmp = append(ptmp, 28)
		}
		pos = append(pos, ptmp)
		ptmp = make([]int, 0)
	}
	verp := make([]float64, times) //儲存roundwin贏分

	binfile, err := os.OpenFile("./payout.bin", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0604)
	if err != nil {
		fmt.Println("開啟失敗", err.Error())
		return
	}
	defer binfile.Close()

	var bin_buf bytes.Buffer
	for i = 0; i < times; i++ {
		loadber(i)
		costmoney = costmoney + len(user.Mg.Paylines) //押注金額
		//產生盤面
		for h = 0; h < len(user.Mg.Field); h++ {
			r := rand.Intn(len(user.Mg.Stripe90[h]))
			//println(r)

			//判斷是否有超過輪帶的長度有則返回0
			for j = 0; j < user.Mg.Field[h]; j++ {
				if r+j > len(user.Mg.Stripe90[h])-1 {
					pos[h][j] = user.Mg.Stripe90[h][r+j-len(user.Mg.Stripe90[h])]

				} else {
					pos[h][j] = user.Mg.Stripe90[h][r+j]
				}
				binary.Write(&bin_buf, binary.LittleEndian, int16(pos[h][j]))
				b := bin_buf.Bytes()
				_, err = binfile.Write(b)

				if err != nil {
					fmt.Println("盤面寫入失敗", err.Error())
					return
				}
			}
		}
		fmt.Println("盤面寫入成功", err.Error())

		//印出盤面
		/*var printstrip [][]string
		var tmp []string
		for j = 0; j < len(user.Mg.Field); j++ {
			for h = 0; h < max(user.Mg.Field); h++ {
				tmp = append(tmp, user.Common.Symbols[pos[j][h]])
			}
			printstrip = append(printstrip, tmp)
			tmp = make([]string, 0)
		}
		transpose(printstrip)*/

		//計算連線分數
		var j, k int
		var mw, mwf, fw int = 0, 0, 0
		for j = 0; j < len(user.Mg.Paylines); j++ {
			//println("line_", j+1, ":", user.Common.Symbols[pos[0][user.Mg.Paylines[j][0]]], user.Common.Symbols[pos[1][user.Mg.Paylines[j][1]]], user.Common.Symbols[pos[2][user.Mg.Paylines[j][2]]], user.Common.Symbols[pos[3][user.Mg.Paylines[j][3]]], user.Common.Symbols[pos[4][user.Mg.Paylines[j][4]]])

			var count, jup int = 0, 0
			if pos[0][user.Mg.Paylines[j][0]] != 15 {
				for k = 0; k < 4; k++ {
					if pos[0][user.Mg.Paylines[j][0]] == pos[k+1][user.Mg.Paylines[j][k+1]] || pos[k+1][user.Mg.Paylines[j][k+1]] == 15 {
						count = count + 1
					} else {
						break
					}
				}
				mw = mw + user.Mg.Paytable[pos[0][user.Mg.Paylines[j][0]]][count]
				//SymbolWin[pos[0][user.Mg.Paylines[j][0]]][count] = SymbolWin[pos[0][user.Mg.Paylines[j][0]]][count] + user.Mg.Paytable[pos[0][user.Mg.Paylines[j][0]]][count]
				//println(user.Common.Symbols[pos[0][user.Mg.Paylines[j][0]]], "*", count+1, "=", user.Mg.Paytable[pos[0][user.Mg.Paylines[j][0]]][count])
			} else if pos[0][user.Mg.Paylines[j][0]] == 15 {
				for k = 0; k < 4; k++ {
					if pos[k+1][user.Mg.Paylines[j][k+1]] == 15 {
						count = count + 1
					} else {
						break
					}
				}
				if count == 4 {
					mw = mw + user.Mg.Paytable[15][4]
					//SymbolWin[15][4] = SymbolWin[15][4] + user.Mg.Paytable[15][4]
					//println(user.Common.Symbols[15], "*", count, "=", user.Mg.Paytable[15][count])
				} else if count == 3 {
					jup = count + 1
					if user.Mg.Paytable[15][count] > user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] {
						mw = mw + user.Mg.Paytable[15][count]
						//SymbolWin[15][count] = SymbolWin[15][count] + user.Mg.Paytable[15][count]
						//println(user.Common.Symbols[15], "*", count+1, "=", user.Mg.Paytable[15][count])
					} else {
						mw = mw + user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup]
						//SymbolWin[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] = SymbolWin[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] + user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup]
						//println(user.Common.Symbols[pos[count+1][user.Mg.Paylines[j][count+1]]], "*", jup+1, "=", user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup])
					}
				} else if count == 2 {
					jup = count + 1
					for k = 3; k < 4; k++ {
						if (pos[count+1][user.Mg.Paylines[j][count+1]] == pos[k+1][user.Mg.Paylines[j][k+1]]) || pos[k+1][user.Mg.Paylines[j][k+1]] == 15 {
							jup = jup + 1
						} else {
							break
						}
					}
					if user.Mg.Paytable[15][count] > user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] {
						mw = mw + user.Mg.Paytable[15][count]
						//SymbolWin[15][count] = SymbolWin[15][count] + user.Mg.Paytable[15][count]
						//println(user.Common.Symbols[15], "*", count+1, "=", user.Mg.Paytable[15][count])
					} else {
						mw = mw + user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup]
						//SymbolWin[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] = SymbolWin[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] + user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup]
						//println(user.Common.Symbols[pos[count+1][user.Mg.Paylines[j][count+1]]], "*", jup+1, "=", user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup])
					}
				} else if count == 1 {
					jup = count + 1
					for k = 2; k < 4; k++ {
						if (pos[count+1][user.Mg.Paylines[j][count+1]] == pos[k+1][user.Mg.Paylines[j][k+1]]) || pos[k+1][user.Mg.Paylines[j][k+1]] == 15 {
							jup = jup + 1
						} else {
							break
						}
					}
					if user.Mg.Paytable[15][count] > user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] {
						mw = mw + user.Mg.Paytable[15][count]
						//SymbolWin[15][count] = SymbolWin[15][count] + user.Mg.Paytable[15][count]
						//println(user.Common.Symbols[15], "*", count+1, "=", user.Mg.Paytable[15][count])
					} else {
						mw = mw + user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup]
						//SymbolWin[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] = SymbolWin[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] + user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup]
						//println(user.Common.Symbols[pos[count+1][user.Mg.Paylines[j][count+1]]], "*", jup+1, "=", user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup])
					}
				} else if count == 0 {
					jup = count + 1
					for k = 1; k < 4; k++ {
						if (pos[count+1][user.Mg.Paylines[j][count+1]] == pos[k+1][user.Mg.Paylines[j][k+1]]) || pos[k+1][user.Mg.Paylines[j][k+1]] == 15 {
							jup = jup + 1
						} else {
							break
						}
					}
					if user.Mg.Paytable[15][count] > user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] {
						mw = mw + user.Mg.Paytable[15][count]
						//SymbolWin[15][count] = SymbolWin[15][count] + user.Mg.Paytable[15][count]
						//println(user.Common.Symbols[15], "*", count+1, "=", user.Mg.Paytable[15][count])
					} else {
						mw = mw + user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup]
						//SymbolWin[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] = SymbolWin[pos[count+1][user.Mg.Paylines[j][count+1]]][jup] + user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup]
						//println(user.Common.Symbols[pos[count+1][user.Mg.Paylines[j][count+1]]], "*", jup+1, "=", user.Mg.Paytable[pos[count+1][user.Mg.Paylines[j][count+1]]][jup])
					}
				}
			}
		}
		mgtotalwin = mgtotalwin + mw
		//Freegame判斷
		var fgsym int //Scatter個數計算
		for h = 0; h < len(user.Mg.Field); h++ {
			for j = 0; j < user.Mg.Field[h]; j++ {
				if pos[h][j] == 18 {
					fgsym = fgsym + 1
				}
			}
		}
		if fgsym == 3 {
			mwf = mwf + user.Mg.Quapaytable[0][2]*len(user.Mg.Paylines)
			//SymbolWin[18][2] = SymbolWin[18][2] + user.Mg.Quapaytable[0][2]*len(user.Mg.Paylines)
			fghit3 = fghit3 + 1
			fw = fgspin(8)
			fgtotalwin3 = fgtotalwin3 + fw
		} else if fgsym == 4 {
			mwf = mwf + user.Mg.Quapaytable[0][3]*len(user.Mg.Paylines)
			//SymbolWin[18][3] = SymbolWin[18][3] + user.Mg.Quapaytable[0][3]*len(user.Mg.Paylines)
			fghit4 = fghit4 + 1
			fw = fgspin(12)
			fgtotalwin4 = fgtotalwin4 + fw
		} else if fgsym == 5 {
			mwf = mwf + user.Mg.Quapaytable[0][4]*len(user.Mg.Paylines)
			//SymbolWin[18][4] = SymbolWin[18][4] + user.Mg.Quapaytable[0][4]*len(user.Mg.Paylines)
			fghit5 = fghit5 + 1
			fw = fgspin(15)
			fgtotalwin5 = fgtotalwin5 + fw
		}

		mgsctotalwin = mgsctotalwin + mwf
		fgtotalwin = fgtotalwin + fw
		mgtotalwin = mgtotalwin + mwf
		//println("===========================================")
		//HITrate WINrate計算
		roundwin = mw + fw + mwf
		verp[i] = float64(roundwin)
		if roundwin > 0 {
			hitp = hitp + 1
		}
		if roundwin > len(user.Mg.Paylines) {
			winp = winp + 1
		}
	}

	//CV值的計算
	var average, sur float64
	average = float64(mgtotalwin+fgtotalwin) / float64(times)
	for i = 0; i < times; i++ {
		sur = sur + math.Pow((verp[i]-average), 2)
	}
	cv = math.Sqrt(sur/float64(times)) / average

	rtp = float64(mgtotalwin+fgtotalwin) / float64(costmoney)
	mgrtp = float64(mgtotalwin) / float64(costmoney)
	fgrtp = float64(fgtotalwin) / float64(costmoney)
	mgScrtp = float64(mgsctotalwin) / float64(costmoney)
	fghit3percent = fghit3 / float64(times)
	fghit4percent = fghit4 / float64(times)
	fghit5percent = fghit5 / float64(times)
	fgretrigerhit3per = fgretrigerhit3 / float64(kkk)
	fgretrigerhit4per = fgretrigerhit4 / float64(kkk)
	fgretrigerhit5per = fgretrigerhit5 / float64(kkk)
	fgrtp3 = float64(fgtotalwin3) / float64(costmoney)
	fgrtp4 = float64(fgtotalwin4) / float64(costmoney)
	fgrtp5 = float64(fgtotalwin5) / float64(costmoney)
	hitpercent = hitp / float64(times)
	winpercent = winp / float64(times)
	for i = 0; i < 11; i++ {
		for j = 0; j < 5; j++ {
			SymbolWin2[i][j] = Round(float64(SymbolWin[i][j])/float64(costmoney)*100, 4)
		}
	}

	binary.Write(&bin_buf, binary.LittleEndian, int16(mgtotalwin+fgtotalwin))
	b := bin_buf.Bytes()
	_, err = binfile.Write(b)

	if err != nil {
		fmt.Println("寫入失敗", err.Error())
		return
	}
	binfile.Close()
}

// FreeGame------------------------------------------------
func fgspin(times int) int {
	//var len(user.Mg.Field) int = len(user.Fg1.Field)
	var i, h int
	var w int = 0
	var posnum int = 1
	for i = 0; i < len(user.Fg1.Field); i++ {
		posnum = posnum * user.Fg1.Field[i]
	}
	pos := make([]int, posnum) //盤面
	for i = 1; i <= times; i++ {
		kkk = kkk + 1
		for h = 0; h < 5; h++ {
			r := rand.Intn(len(user.Fg1.Stripe90[h]))
			pos[3*h] = user.Fg1.Stripe90[h][r]
			if r == len(user.Fg1.Stripe90[h])-1 {
				pos[3*h+1] = user.Fg1.Stripe90[h][0]
				pos[3*h+2] = user.Fg1.Stripe90[h][1]
			} else if r == len(user.Fg1.Stripe90[h])-2 {
				pos[3*h+1] = user.Fg1.Stripe90[h][r+1]
				pos[3*h+2] = user.Fg1.Stripe90[h][0]
			} else {
				pos[3*h+1] = user.Fg1.Stripe90[h][r+1]
				pos[3*h+2] = user.Fg1.Stripe90[h][r+2]
			}
		}
		var fgretriger int
		for h = 0; h < 15; h++ {
			if pos[h] == 18 {
				fgretriger = fgretriger + 1
			}
		}
		//Freegame判斷
		if fgretriger == 3 {
			w = w + user.Fg1.Quapaytable[0][2]*len(user.Fg1.Paylines)
			if times+8 <= 40 {
				times = times + 8
			} else {
				times = 40
			}
			fgretrigerhit3 = fgretrigerhit3 + 1
		} else if fgretriger == 4 {
			w = w + user.Fg1.Quapaytable[0][3]*len(user.Fg1.Paylines)
			if times+12 <= 40 {
				times = times + 12
			} else {
				times = 40
			}
			fgretrigerhit4 = fgretrigerhit4 + 1
		} else if fgretriger == 5 {
			w = w + user.Fg1.Quapaytable[0][4]*len(user.Fg1.Paylines)
			if times+15 <= 40 {
				times = times + 15
			} else {
				times = 40
			}
			fgretrigerhit5 = fgretrigerhit5 + 1
		}

		//計算連線分數
		var j, k int
		for j = 0; j < len(user.Fg1.Paylines); j++ {
			var count int = 0
			if pos[user.Fg1.Paylines[j][0]] != 15 {
				for k = 0; k < 4; k++ {
					if pos[user.Fg1.Paylines[j][0]] == pos[(k+1)*3+user.Fg1.Paylines[j][k+1]] || pos[(k+1)*3+user.Fg1.Paylines[j][k+1]] == 15 {
						count = count + 1
					} else {
						break
					}
				}
				w = w + user.Fg1.Paytable[pos[user.Fg1.Paylines[j][0]]][count]
			} else if pos[user.Fg1.Paylines[j][0]] == 15 {
				for k = 0; k < 4; k++ {
					if pos[(k+1)*3+user.Fg1.Paylines[j][k+1]] == 15 {
						count = count + 1
					} else {
						break
					}
				}
				if count == 4 {
					w = w + user.Fg1.Paytable[15][4]
				} else if count == 3 {
					if user.Fg1.Paytable[15][3] > user.Fg1.Paytable[pos[4*3+user.Fg1.Paylines[j][4]]][4] {
						w = w + user.Fg1.Paytable[15][3]
					} else {
						w = w + user.Fg1.Paytable[pos[4*3+user.Fg1.Paylines[j][4]]][4]
					}
				} else if count == 2 {
					count = count + 1
					for k = 3; k < 4; k++ {
						if (pos[3*3+user.Fg1.Paylines[j][3]] == pos[(k+1)*3+user.Fg1.Paylines[j][k+1]]) || pos[(k+1)*3+user.Fg1.Paylines[j][k+1]] == 15 {
							count = count + 1
						} else {
							break
						}
					}
					if user.Fg1.Paytable[15][2] > user.Fg1.Paytable[pos[3*3+user.Fg1.Paylines[j][3]]][count] {
						w = w + user.Fg1.Paytable[15][2]
					} else {
						w = w + user.Fg1.Paytable[pos[3*3+user.Fg1.Paylines[j][3]]][count]
					}
				} else if count == 1 {
					count = count + 1
					for k = 2; k < 4; k++ {
						if (pos[2*3+user.Fg1.Paylines[j][2]] == pos[(k+1)*3+user.Fg1.Paylines[j][k+1]]) || pos[(k+1)*3+user.Fg1.Paylines[j][k+1]] == 15 {
							count = count + 1
						} else {
							break
						}
					}
					if user.Fg1.Paytable[15][1] > user.Fg1.Paytable[pos[2*3+user.Fg1.Paylines[j][2]]][count] {
						w = w + user.Fg1.Paytable[15][1]
					} else {
						w = w + user.Fg1.Paytable[pos[2*3+user.Fg1.Paylines[j][2]]][count]
					}
				} else if count == 0 {
					count = count + 1
					for k = 1; k < 4; k++ {
						if (pos[3+user.Fg1.Paylines[j][1]] == pos[(k+1)*3+user.Fg1.Paylines[j][k+1]]) || pos[(k+1)*3+user.Fg1.Paylines[j][k+1]] == 15 {
							count = count + 1
						} else {
							break
						}
					}
					if user.Fg1.Paytable[15][0] > user.Fg1.Paytable[pos[3+user.Fg1.Paylines[j][1]]][count] {
						w = w + user.Fg1.Paytable[15][0]
					} else {
						w = w + user.Fg1.Paytable[pos[3+user.Fg1.Paylines[j][1]]][count]
					}
				}
			}
		}
	}
	return w
}
