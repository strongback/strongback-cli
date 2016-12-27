## Strongback Command Line Interface (CLI)

[Strongback](http://strongback.org) is a Java library for FRC robots, and to use it you must first install the library onto your development machine and set it up properly. This can be tedious and somewhat complicated, so the Strongback Command Line Interface (CLI) tool provides an easy-to-use way to interact with the library using the command line. The Strongback CLI consists of a single, entirely self-contained executable program -- simply download the right program for your operating system and then run it to:

* list the Strongback Java Library versions that are available
* install, upgrade, restore, or uninstall the Strongback Java Library version
* create a new Java project that uses Strongback
* decode a binary data file recorded by the Strongback library running on a robot
* show information about the installed Strongback and WPILib for Java

## Installing the Strongback CLI

There are separate Strongback Command Line Interface (CLI) executables for these operating systems:

* Windows 32-bit (x86)
* Linux (64-bit x86)
* Mac OS X 10.7 and higher

The instructions for Windows differ slightly from those for OS X and Linux.

### Windows

Look at the [latest releases](https://github.com/strongback/strongback-cli/releases) of this utility, and download the file that ends with `-windows.zip` or `-windows.tar.gz` for the latest stable release. The archive contains a single `<strongback.exe` file, so unzip or unpack the archive into any directory that is already included in or that [you can add to the `%PATH%` environment variable](http://www.howtogeek.com/118594/how-to-edit-your-system-path-for-easy-command-line-access/). We recommend your home directory, e.g., `C:\\Users\\<you>`.

After you've extracted the file and [added the directory to your `%PATH%`](http://www.howtogeek.com/118594/how-to-edit-your-system-path-for-easy-command-line-access/), open up a command window and run the following command from several different directories:

    > strongback info

This will show you the installed versions of the Strongback Java Library and WPILib for Java library.

### Linux and OS X

Simply download the [latest version](https://github.com/strongback/strongback-cli/releases) of the `tar.gz` file and unpack it into any directory that is already on your `PATH` _except your home directory_. We recommend a `bin` directory in your home account.

To create the `bin` directory, use the following commands:

    $ cd ~
    $ mkdir -p ~/bin

which will work even if the `~/bin` directory exists. Add the `~/bin` directory to your path:

    $ echo "PATH=\${HOME}/bin:\${PATH}" >> ~/.bashrc
    $ source ~/.bashrc

Then, download and install the `strongback` executable into your `bin` directory. You can use the following commands, though be sure to replace both "1.0.0" sequences with the [latest version](https://github.com/strongback/strongback-cli/releases):

For OS X:

    $ cd ~/bin
    $ curl -o https://github.com/strongback/strongback-cli/releases/download/v1.0.0/strongback-cli-1.0.0-osx.tar.gz | tar xvz

For Linux:

    $ cd ~/bin
    $ curl -o https://github.com/strongback/strongback-cli/releases/download/v1.0.0/strongback-cli-1.0.0-linux.tar.gz | tar xvz

That's it! Open up a new terminal and run the following:

    $ strongback info

This will show you the installed versions of the Strongback Java Library and WPILib for Java library.

## Viewing help

The Strongback CLI has built-in help, which you can see by running `strongback help` to display something like:

    Usage:

       strongback <command> [<args>]

    Available commands include:
       decode        Converts a binary data/event log file to a readable CSV file
       help          Displays information about using this utility
       info          Displays the information about this utility and what's installed
       install       Install or upgrade the Strongback Java Library
       new-project   Creates a new project configured to use Strongback (only 1.x)
       releases      Display the available versions of the Strongback Java Library
       version       Display the currently installed version
       uninstall     Remove an existing Strongback Java Library installation

    Additional help is available for each command with:

       strongback help <command>

We've already seen the `info` command, and the `version` is similar but more concise. Let's look at several other commands.

## Installing or upgrading Strongback Java Library

The Strongback CLI makes it easy to list and install the available versions of the Strongback Java Library and to install any of these. Run the following command to list the available versions:

    $ strongback releases

This will check the [Strongback releases](https://github.com/strongback/strongback-cli/releases) and output something similar to:

    Found 10 releases of the Strongback Java Library:
      1.1.7
      1.1.6
      1.1.5
      1.1.3
      1.1.2
      1.1.1
      1.1.0
      1.0.3
      1.0.2
      1.0.1

To install one of these releases, simply run `strongback install <version>` and supply the desired version number (e.g., `strongback install 1.1.7`). Or, if you want the latest version, simply run `strongback install`. If you already have that version installed, the tool will simply tell you this and return. Otherwise, it will archive your existing version (if you have one installed) and then install the version you specified.

You can use this same command to restore a version that was installed previously, allowing you to easily switch between different installed versions. For example, imagine that you've recently installed version 1.1.6 but want to try 1.1.7. You can upgrade to 1.1.7 with `strongback install 1.1.7` and then later switch back to your previous 1.1.6 installation with `strongback install 1.1.6`. To use 1.1.7, simply run `strongback install 1.1.7` again.

Whenever you install a new version, the CLI will archive any previously installed version in the `~/strongback-archives` directory. Reinstalling one of these simply extracts that archive rather than downloading the archive from the [Strongback releases](https://github.com/strongback/strongback-cli/releases).

## Creating a new robot project

You can use the Strongback CLI tool to create a new iterative Java robot project for Eclipse set up to use Strongback, or update an existing Java robot project to use Strongback. Simply run open a terminal, change to the directory where you want the project created, and run the following command:

    $ strongback new-project MyRobotProject

and replace `MyRobotProject` with the name of your project. This will not overwrite any of the existing files, so to do that add the `--overwrite` flag:

    $ strongback new-project --overwrite MyRobotProject

By default, the Java package will be `org.frc<teamNumber>` where `<teamNumber>` is read from your WPILib for Java installation. If you want to use a different package, then supply the `--package <packageName>` flag. For example:

    $ strongback new-project --package io.alphateam.robot MyRobotProject

See `strongback help new-project` for additional options.

## Decoding a robot data log

The Strongback Java Library has a _data recorder_ capable of recording various channels of data while your robot runs. The resulting data is captured in a binary file on the robot, so you need to download the file from the robot and then decode it. The Strongback CLI's `decode` command will convert the binary file into a _comma separated values_ (or CSV) file that you can import into Google Sheets, Excel, Tableau, or many other programs.

For example, imagine that you've set up your robot to use Strongback's data recorder to capture two channels, and you've named these channels `Foo` and `Bar`. The raw binary file will look something like the following (new lines and brackets added for clarity and are not part of the file format):

    [l o g]
    [3]
    [4] [2] [2]
    [4][T i m e] [3][F o o] [3][B a r]
    [00 00 00 00] [00 52] [00 37]
    [00 00 00 0A] [04 D5] [23 AF]
    [00 00 00 14] [3F 00] [12 34]
    [FF FF FF FF]

If this is in a file named `robot.dat` downloaded from the robot onto your computer, then the following Strongback CLI command will convert this binary file to a CSV file:

    $ strongback decode robot.dat robot.csv

The CSV file will look like:

    Time, Foo, Bar
    0, 82, 55,
    10, 1237, 9135,
    20, 16128, 4660

The first row contains the comma-separated names of the channels, followed by a newline character. Each of the subsequent line lists the integer value of the data encoded to the precision specified in your robot program. The `Time` channel is always first and the values are in milliseconds.


## Uninstalling the Strongback Java Library

To uninstall the Strongback Java Library, simply run

    $ strongback uninstall

and your existing installation (if you have one) will be archived.

## Uninstalling all of Strongback

To fully uninstall all of Strongback, including the archives of existing installations, run the following commands:

    $ strongback uninstall --remove-archives

and then remove the `strongback` CLI tool.

