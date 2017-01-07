# coco

Command line tool to count single lines of code (SLOC).

Similar to the classic wc command, but counts comments and empty lines as
well. Lines that contain both code and comments are counted only as code.

## Usage

    Usage: coco [options] [files...]
    Prints a summary of total, non-commented, commented and empty lines.

    The default comment style is C. To change it, use the -s and -m options.

    If no files are given, standard input is used.

    Example:
        coco -s "#" -m "" main.py
        cat lib.c | coco

    -h
    -help
            Show usage
    -m string
            Alias for -multi (default "/* */")
    -multi string
            Multi-line comment with start and end separated by a space. (default "/* */")
    -s string
            Alias for -single (default "//")
    -single string
            Single-line comment (default "//")
    -v
    -verbose
            Output title for each count

## License

Free. As in free beer.
