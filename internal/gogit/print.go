package gogit

import "fmt"

// PrintCommit prints a commit object with a stylized format.
func PrintCommit(commit *Commit) {
	fmt.Printf("%scommit %s%s\n", ColorYellow, commit.Hash, ColorReset)
	fmt.Printf("Tree: %s\n", commit.Tree)
	if commit.Parent != "" {
		fmt.Printf("%sParent: %s%s\n", ColorRed, commit.Parent, ColorReset)
	}
	fmt.Printf("%sAuthor: %s%s\n", ColorGreen, commit.Author, ColorReset)
	fmt.Printf("%sDate: %s%s\n", ColorBlue, commit.Date.Format("Mon Jan 2 15:04:05 2006 -0700"), ColorReset)
	fmt.Printf("\n\t%s\n\n", commit.Message)
}

func PrintStatus(statusInfo *StatusInfo) {
	// Imprimir la rama actual
	fmt.Printf("On branch %s\n", statusInfo.Branch)

	// Variable para saber si el repositorio está limpio
	isClean := true

	// Mostrar archivos listos para commit (Staged)
	if len(statusInfo.Staged) > 0 {
		isClean = false
		fmt.Println("\nChanges to be committed:")
		fmt.Println("  (use \"go-git reset <file>...\" to unstage)") // Cambia "go-git" por el nombre de tu app
		for _, file := range statusInfo.Staged {
			fmt.Printf("%s\t%s%s\n", ColorGreen, file, ColorReset)
		}
	}

	// Mostrar archivos con cambios no preparados para commit (Unstaged)
	if len(statusInfo.Unstaged) > 0 {
		isClean = false
		fmt.Println("\nChanges not staged for commit:")
		fmt.Println("  (use \"go-git add <file>...\" to update what will be committed)") // Cambia "go-git" por el nombre de tu app
		for _, file := range statusInfo.Unstaged {
			fmt.Printf("%s\t%s%s\n", ColorRed, file, ColorReset)
		}
	}

	// Mostrar archivos no seguidos (Untracked)
	if len(statusInfo.Untracked) > 0 {
		isClean = false
		fmt.Println("\nUntracked files:")
		fmt.Println("  (use \"go-git add <file>...\" to include in what will be committed)") // Cambia "go-git" por el nombre de tu app
		for _, file := range statusInfo.Untracked {
			fmt.Printf("%s        %s%s\n", ColorRed, file, ColorReset)
		}
	}

	// Si no hubo cambios en ninguna sección, el árbol de trabajo está limpio
	if isClean {
		fmt.Println("\nnothing to commit, working tree clean")
	}
}
