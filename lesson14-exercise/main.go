package main

import (
	"fmt"

	"galvin/lession14-exercise/student"
	"galvin/lession14-exercise/teacher"
	"galvin/lession14-exercise/utils"
)

func main() {
	
	for {
		utils.ClearScreen()

		fmt.Println("📚 CHUONG TRINH QUAN LY")
		fmt.Println("1. Quan ly sinh vien")
		fmt.Println("2. Quan ly giang vien")
		fmt.Println("3. Thoat")

		choice := utils.GetPostiveInt("👉 Chon chuc nang: ")

		switch choice {
		case 1:
			student.StudentMenu()
		case 2:
			teacher.TeacherMenu()
		case 3:
			return
		default:
			fmt.Println("❌ Lua chon khong hop le!")
		}
	
		utils.ReadInput("\nNhan phim Enter de tiep tuc...")
	}
	

}
