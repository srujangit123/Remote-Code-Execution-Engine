#!/bin/sh

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <language> <source_file_path> <output_file_path>"
    exit 1
fi

language="$1"
source_file="$2"
output_file="$3"

if [ ! -f "$source_file" ]; then
    echo "Error: Source file '$source_file' does not exist!" >> "$output_file"
    exit 1
fi

run_cpp() {
    executable="a.out"

    g++ "$source_file" -o "$executable" 2> "$output_file"_compile_errors.txt

    if [ $? -ne 0 ]; then
        echo "Compilation failed. Check 'compile_errors.txt' for details." >> "$output_file"
        cat "$output_file"_compile_errors.txt >> "$output_file"
        rm compile_errors.txt
        exit 1
    fi

    ./$executable > "$output_file" 2>&1

    if [ $? -ne 0 ]; then
        echo "Runtime error occurred. Check the output file for details." >> "$output_file"
    fi

    rm "$executable"
}

run_go() {
    # First, compile the Go program to check for compile-time errors
    # go build "$source_file" -o "$output_file"_exe 2> "$output_file"
    go run "$source_file" > "$output_file" 2>&1

    # Check if compilation failed
    # if [ $? -ne 0 ]; then
    #     echo "Go compilation failed. Check $output_file" >> "$output_file"
    #     exit 1
    # fi

    # Now, run the Go program if compilation is successful
    # ./"$output_file"_exe > "$output_file" 2>&1

    # Check for runtime errors
    if [ $? -ne 0 ]; then
        echo "Runtime error occurred while running the Go program." >> "$output_file"
    fi

    # Clean up the compiled executable
    # rm -f "$output_file"_exe
}

# Check the programming language and call the appropriate function
case "$language" in
    cpp)
        run_cpp
        ;;
    golang)
        run_go
        ;;
    *)
        echo "Error: Unsupported language '$language'. Please use 'cpp' or 'go'." >> "$output_file"
        exit 1
        ;;
esac
