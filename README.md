# GoGit - A Git Replica in Go

GoGit is a simplified implementation of Git, the distributed version control system, built from scratch in Go. This project is intended for educational purposes, providing a hands-on approach to understanding the core concepts and inner workings of Git.

## Features

*   **Initialize a new repository:** Create a new `.gogit` directory to start tracking a project.
*   **Add files to the staging area:** Track new or modified files to be included in the next commit.
*   **Commit changes:** Save snapshots of the staging area to the project's history.
*   **View commit history:** Inspect the log of commits to see the project's evolution.

## Getting Started

### Prerequisites

*   Go 1.18 or higher

### Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/TonyGLL/go-git.git
    ```
2.  Navigate to the project directory:
    ```sh
    cd go-git
    ```
3.  Build the project:
    ```sh
    go build .
    ```

## Usage

GoGit provides the following commands:

*   `gogit init`: Initializes a new repository.
*   `gogit add <file>`: Adds a file to the staging area.
*   `gogit commit -m <message>`: Commits the staged changes.
*   `gogit log`: Displays the commit history.

## Contributing

Contributions are welcome! If you'd like to improve GoGit, please feel free to fork the repository, make your changes, and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
