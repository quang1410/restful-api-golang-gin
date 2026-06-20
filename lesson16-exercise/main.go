package main

import (
	"fmt"

	"galvin/lesson16-exercise/library"
	"galvin/lesson16-exercise/utils"
)

func main() {

	lib := library.NewLibrary()

	for {
		utils.ClearScreen()

		fmt.Println("📚 CHUONG TRINH QUAN LY THU VIEN")
		fmt.Println("1. Them sach")
		fmt.Println("2. Xem danh sach sach")
		fmt.Println("3. Them nguoi muon")
		fmt.Println("4. Xem danh sach nguoi muon")
		fmt.Println("5. Muon sach")
		fmt.Println("6. Xem lich su muon")
		fmt.Println("7. Tra sach")
		fmt.Println("8. Tim kiem sach")
		fmt.Println("9. Thoat")

		choice := utils.GetPostiveInt("👉 Chon chuc nang: ")

		utils.ClearScreen()

		switch choice {
		case 1:
			fmt.Println("-=-=-=-=- Them Sach -=-=-=-=-")
			if err := library.AddBook(lib); err != nil {
				fmt.Printf("❌ Loi khi them sach: %v\n", err)
			}
		case 2:
			fmt.Println("-=-=-=-=- Xem Danh Sach Sach -=-=-=-=-")
			if err := library.ListBooks(lib); err != nil {
				fmt.Printf("❌ Loi khi xem danh sach sach: %v\n", err)
			}
		case 3:
			fmt.Println("-=-=-=-=- Them Nguoi Muon Sach -=-=-=-=-")
			if err := library.AddBorrower(lib); err != nil {
				fmt.Printf("❌ Loi khi them nguoi muon sach: %v\n", err)
			}
		case 4:
			fmt.Println("-=-=-=-=- Xem Danh Sach Nguoi Muon -=-=-=-=-")
			if err := library.ListBorrowers(lib); err != nil {
				fmt.Printf("❌ Loi khi xem danh sach nguoi muon: %v\n", err)
			}
		case 5:
			fmt.Println("-=-=-=-=- Muon Sach -=-=-=-=-")
			if err := library.BorrowBook(lib); err != nil {
				fmt.Printf("❌ Loi khi muon sach: %v\n", err)
			}
		
		case 6:
			fmt.Println("-=-=-=-=- Xem Lich Su Muon Sach -=-=-=-=-")
			if err := library.ListBorrowHistory(lib); err != nil {
				fmt.Printf("❌ Loi khi xem lich su muon sach: %v\n", err)
			}
		case 7:
			fmt.Println("-=-=-=-=- Tra Sach -=-=-=-=-")
			if err := library.ReturnBook(lib); err != nil {
				fmt.Printf("❌ Loi khi tra sach: %v\n", err)
			}
		case 8:
			fmt.Println("-=-=-=-=- Tim Kiem Sach -=-=-=-=-")
			if err := library.SearchBooks(lib); err != nil {
				fmt.Printf("❌ Loi khi tim kiem sach: %v\n", err)
			}
		case 9:
			fmt.Println("Tam biet!")
			return
		default:
			fmt.Println("❌ Lua chon khong hop le!")
		}

		utils.ReadInput("\nNhan Enter de tiep tuc ")

	}	

}