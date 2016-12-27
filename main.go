package main

import (
    "net/http"
    "bufio"
    "bytes"
    "time"
    "strings"
    "encoding/json"
    "fmt"
    "flag"
    "io"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "github.com/rickar/props"
    "strongback.org/cli/files"
)

var (
    Version   string
    Date      string
    Build     string
    ExecName  string
)

type Environment struct {
    userHome string
    existingStrongback Component
    existingWpiLib Component
    httpClient http.Client
    availableReleases []ReleaseInfo
    teamNumber string
}

type Component struct {
    installed bool
    version string
    path string
    properties *props.Properties
}

type ReleaseInfo struct {
    Name string
    Url string
    Html_url string
    Id int
    Tag_name string
    Draft bool
    Prerelease bool
    Published_at string
    Assets []AssetInfo
}

type AssetInfo struct {
    Url string
    Browser_download_url string
    Id int
    Name string
    Label string
    Content_type string
    Size int
}

// NewProperties creates a new environment.
func NewEnvironment() *Environment {
    e := new(Environment)
    e.httpClient = http.Client{Timeout: 10 * time.Second}
    e.userHome = files.UserHomeDir()
    e.teamNumber = ""
    e.DiscoverStrongback()
    e.DiscoverWpiLib()
    return e
}

func (env *Environment) DiscoverStrongback() {
    dir := env.userHome + files.PathSeparator + "strongback"
    strongback := new(Component)
    strongback.path = dir
    strongback.version = "<none>"
    strongback.installed = false
    if files.IsExistingDirectory(dir) {
        // Load the properties file ...
        propPath := dir + files.PathSeparator + "strongback.properties"
        if files.IsExistingFile(propPath) {
            props := *files.LoadPropertiesFile(propPath)

            // Create the installed component ...
            strongback.properties = &props
            if &props != nil && props.Get("strongback.version") != "" {
                strongback.version = props.Get("strongback.version")
                strongback.installed = true
            }
        }
    }

    env.existingStrongback = *strongback
}

func (env *Environment) DiscoverWpiLib() {
    dir := env.userHome + files.PathSeparator + "wpilib"
    propPath := dir + files.PathSeparator + "wpilib.properties"

    props := *files.LoadPropertiesFile(propPath)

    // Create the installed component ...
    wpilib := new(Component)
    wpilib.path = dir
    wpilib.properties = &props
    if &props != nil {
        wpilib.version = props.Get("version")
        wpilib.installed = true
        teamNumberStr := props.Get("team-number")
        if len(teamNumberStr) != 0 {
            env.teamNumber = teamNumberStr
        }
    }

    env.existingWpiLib = *wpilib
}

func (env *Environment) Print() {
    fmt.Println("Strongback Library")
    if env.existingStrongback.installed {
        fmt.Println("  version:      " + env.existingStrongback.version)
        fmt.Println("  build date:   " + env.existingStrongback.properties.Get("build.date"))
        fmt.Println("  location:     " + env.existingStrongback.path)
    } else {
        fmt.Println("  not yet installed (use 'strongback install' command)")
    }
    if env.existingWpiLib.installed {
        fmt.Println("WPILib")
        //fmt.Println("  version:    " + env.existingWpiLib.version)
        fmt.Println("  location:   " + env.existingWpiLib.path)
    } else {
        fmt.Println("WPILib not installed")
    }
}

func (env *Environment) GetAvailableReleases() []ReleaseInfo {
    if env.availableReleases == nil {
        // Get the available releases ...
        var releases []ReleaseInfo
        env.getJson("https://api.github.com/repos/strongback/strongback-java/releases", &releases)
        env.availableReleases = releases
    }
    return env.availableReleases

}

func (env *Environment) GetLatestRelease(includePreReleases bool) *ReleaseInfo {
    releases := env.GetAvailableReleases()
    for _, release := range releases {
        if !includePreReleases && release.IsPreRelease() {
            // Skip the pre-releases 
            continue
        }
        return &release
    }
    return nil
}

func (env *Environment) GetRelease(version string) *ReleaseInfo {
    if len(version) == 0 {
        return env.GetLatestRelease(false)
    }
    releases := env.GetAvailableReleases()
    for _, release := range releases {
        if release.Name == version {
            return &release
        }
    }
    return nil
}

func (release *ReleaseInfo) IsPreRelease() bool {
    return strings.Contains(release.Name,"Alpha") || strings.Contains(release.Name,"Beta")|| strings.HasPrefix(release.Name,"v")
}

func (env *Environment) PrintReleases(includePreReleases bool) {
    // Get the available releases ...
    releases := env.GetAvailableReleases()
    visibleReleases := 0
    for _, release := range releases {
        if !includePreReleases && release.IsPreRelease() {
            // Skip the pre-releases 
        } else {
            visibleReleases = visibleReleases+1
        }
    }
    fmt.Printf("\nFound %d releases of the Strongback Java Library:\n", visibleReleases)
    currentVersion := env.existingStrongback.version
    for _, release := range releases {
        if release.Name == currentVersion {
            fmt.Printf("  %s (installed)\n",release.Name)
        } else {
            if !includePreReleases && release.IsPreRelease() {
                // Skip the pre-releases 
            } else {
                fmt.Printf("  %s\n",release.Name)
            }
        }
    }
}

func (env *Environment) PrintVersion() {
    if env.existingStrongback.installed {
        fmt.Println("strongback library version " + env.existingStrongback.version)
    } else {
        fmt.Println("strongback library version <none>")
    }
    fmt.Println("strongback cli version " + Version)
}

func (env *Environment) PrintInfo() {
    fmt.Println()
    fmt.Println("Strongback Command Line Interface (CLI)")
    fmt.Println("  version:          " + Version)
    fmt.Println("  build date:       " + Date)
    fmt.Println()
    fmt.Println("Strongback Java Library")
    if env.existingStrongback.installed {
        fmt.Println("  current version:  " + env.existingStrongback.version)
        fmt.Println("  build date:       " + env.existingStrongback.properties.Get("build.date"))
        fmt.Println("  location:         " + env.existingStrongback.path)
    } else {
        fmt.Println("  not yet installed (use 'install' command)")
    }
    fmt.Println()
    fmt.Println("WPILib Java Library")
    if env.existingWpiLib.installed {
        fmt.Println("  location:         " + env.existingWpiLib.path)
        if len(env.teamNumber) != 0 {
            fmt.Println("  team number:      " + env.teamNumber)
        } else {
            fmt.Println("  team number:      <create robot project in Eclipse>")
        }
    } else {
        fmt.Println("  not yet installed")
    }
    fmt.Println()
}

func (env *Environment) getJson(url string, target interface{}) {
    r, err := env.httpClient.Get(url)
    if err != nil {
        panic(err)
    }
    defer r.Body.Close()
    decodeErr := json.NewDecoder(r.Body).Decode(target)
    if decodeErr != nil {
        fmt.Printf("%T\n%s\n%#v\n",decodeErr, decodeErr, decodeErr)
        fmt.Print(decodeErr)
        switch v := decodeErr.(type){
            case *json.SyntaxError:
                fmt.Print(v)
                // fmt.Println(string(body[v.Offset-40:v.Offset]))
        }
        //panic(decodeErr)
    }
}

func (env *Environment) GetUserConfirmation(maxAskTimes int) bool {
    reader := bufio.NewReader(os.Stdin)
    // Confirm removal ...
    for times := 1;  times<=maxAskTimes; times++ {
        fmt.Printf("Are you sure? [y/n] ")
        response, err := reader.ReadString('\n')
        if err != nil {
            return false
        }

        response = strings.ToLower(strings.TrimSpace(response))
        if response == "y" || response == "yes" {
            return true
        } else if response == "n" || response == "no" {
            return false
        }
    }
    return false
}

func (env *Environment) UninstallRelease(skipPrompt bool, skipArchive bool, verbose bool, removeArchive bool) bool {
    fmt.Println()
    if env.existingStrongback.installed {
        fmt.Printf("Removing Java Library version %s\n", env.existingStrongback.version)
        if skipPrompt || env.GetUserConfirmation(3) {
            env.RemoveInstalledRelease(skipArchive, verbose)
            if removeArchive {
                fmt.Println()
                fmt.Println("Removing all archives of previous installations. This cannot be undone!")
                fmt.Println()
                if skipPrompt || env.GetUserConfirmation(3) {
                    archiveDirPath := env.existingStrongback.path + "-archives"
                    os.RemoveAll(archiveDirPath)
                    fmt.Println("Archives of previous installations removed.")
                }
            }
        } else {
            fmt.Println()
            fmt.Println("Existing without uninstalling.")
            os.Exit(1)
        }
        return true
    }
    fmt.Println("No Strongback Java Library version is installed.")
    return false;
}

func (env *Environment) RemoveInstalledRelease(skipArchive bool, verbose bool) bool {
    if env.existingStrongback.installed {
        if !skipArchive {
            // There is an existing release, so archive it
            historyArchiveName := "strongback-" + env.existingStrongback.version + ".tar.gz"
            // First, make sure the archive directory exists
            archiveDirPath := env.existingStrongback.path + "-archives"
            files.MkDir(archiveDirPath)
            // Create a tar.gz file with the existing installation, overwriting any existing archive
            historyArchivePath := archiveDirPath + files.PathSeparator + historyArchiveName
            fmt.Println("   archiving current " + env.existingStrongback.version + " installation to " + historyArchivePath)
            err := files.CreateTar(historyArchivePath, env.userHome, "strongback", false)
            if err != nil {
                panic(err)
            }
        }
        os.RemoveAll(env.existingStrongback.path)
        return true
    }
    return false
}

func (env *Environment) InstallRelease(desiredVersion string, skipArchive bool, verbose bool) bool {
    fmt.Println()
    // See what is already installed
    existingVersion := env.existingStrongback.version
    if desiredVersion == existingVersion && desiredVersion != "" {
        fmt.Printf("%s is already installed\n", desiredVersion)
        return true
    }

    // Get the information for the desired release
    var release *ReleaseInfo
    latestAvailable := ""
    if len(desiredVersion) == 0 {
        release = env.GetLatestRelease(false)
        latestAvailable = "latest available"
    } else {
        release = env.GetRelease(desiredVersion)
    }


    if release == nil {
        fmt.Println("Unable to find and install Strongback Java Library version " + desiredVersion)
        return false
    }

    // Find the asset we want to download
    var desiredAsset *AssetInfo
    for _, asset := range release.Assets {
        if asset.Content_type == "application/x-gzip" {
            desiredAsset = &asset
            break
        }
    }
    if desiredAsset == nil {
        fmt.Println("Unable to find a TAR archive for version " + desiredVersion)
        return false
    }

    fmt.Println("Installing " + latestAvailable + " Strongback Java Library " + release.Name)
    fmt.Println()

    archiveName := desiredAsset.Name

    // See if the asset exists in our archive 
    archiveAssetPath := env.existingStrongback.path + "-archives" + files.PathSeparator + archiveName
    if files.IsExistingFile(archiveAssetPath) {
        fmt.Println("   found previously installed archive at " + archiveAssetPath)
    } else {
        // Download the desired release to a local file
        fmt.Print("   downloading " + archiveName + " to " + archiveAssetPath)
        // file does not exist, so download it
        resp, err := env.httpClient.Get(desiredAsset.Browser_download_url)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()

        // Create the file
        file, err := os.Create(archiveAssetPath)
        if err != nil {
            fmt.Println(err)
            panic(err)
        }
        defer file.Close()

        // Copy the downloaded content into the file
        size, err := io.Copy(file, resp.Body)
        if err != nil {
            panic(err)
        }
        fmt.Printf(" (%v bytes)\n", size)
    }

    if env.existingStrongback.installed {
        if env.RemoveInstalledRelease(skipArchive, verbose) {
            fmt.Println("   replacing existing " + env.existingStrongback.version + " installation with " + desiredVersion + " at " + env.existingStrongback.path)
        }
    } else {
        fmt.Println("   installing at " + env.userHome + files.PathSeparator + "strongback")
    }

    // Install this release
    err := files.ExtractTar(archiveAssetPath, env.userHome, verbose)
    if err != nil {
        panic(err)
    }

    // Update the one we know about
    env.DiscoverStrongback()

    // If there is no Eclipse directory in the Strongback installation, then make it ...
    if !files.IsExistingDirectory(env.existingStrongback.path + filepath.FromSlash("/java/eclipse")) {
        projectName := "initialeclipseproject"
        projectDirPath := env.existingStrongback.path + files.PathSeparator + projectName
        os.RemoveAll(projectDirPath)
        env.NewProject(projectName, env.existingStrongback.path, "", true, true, true)
        os.RemoveAll(projectDirPath)
    }

    return true
}

func (env *Environment) CheckInstalled() {
    if !env.existingStrongback.installed {
        fmt.Println("You must first install the Strongback Java Library using:")
        fmt.Println()
        PrintInstallUsage()
        os.Exit(2)
    }
}

func (env *Environment) DecodeFile(inputFile string, outputFile string, verbose bool) {
    suffix := ".sh"
    if runtime.GOOS == "windows" {
        suffix = ".bat"
    }
    var args []string
    args = append(args, "log-decoder")
    args = append(args, "-f")
    args = append(args, inputFile)
    if len(outputFile) != 0 {
        args = append(args, "-o")
        args = append(args, outputFile)
    }
    commandPath := env.existingStrongback.path + filepath.FromSlash("/java/bin/strongback") + suffix;
    out, err := exec.Command(commandPath, args...).Output()
    if err != nil {
        log.Fatal("Error running " + commandPath + " " + strings.Join(args, " "))
        panic(err)
    }
    fmt.Println(string(out))
}

func (env *Environment) NewProject(name string, directory string, packageName string, eclipse bool, overwrite bool, silent bool) bool {
    suffix := ".sh"
    if runtime.GOOS == "windows" {
        suffix = ".bat"
    }
    if len(packageName) == 0 {
        if len(env.teamNumber) != 0 {
            packageName = "org.frc" + env.teamNumber + ".robot"      
        } else {
            if !silent {
                fmt.Println()
                fmt.Println("No package name was specified, and WPILib has not been initialized with a team number.") 
                fmt.Println("Aborting.")
            }
            return false           
        }
    }
    var args []string
    args = append(args, "new-project")
    args = append(args, "-n")
    args = append(args, name)
    args = append(args, "-d")
    args = append(args, directory)
    args = append(args, "-p")
    args = append(args, packageName)
    if eclipse {
        args = append(args, "-e")
    }
    if overwrite {
        args = append(args, "-o")
    }
    commandPath := env.existingStrongback.path + filepath.FromSlash("/java/bin/strongback") + suffix;
    out, err := exec.Command(commandPath, args...).Output()
    if err != nil {
        log.Fatal("Error running " + commandPath + " " + strings.Join(args, " "))
        panic(err)
    }
    if !silent {
        fmt.Println(string(out))
    }
    return true
}

func PrintUsage() {
    fmt.Println("   " + ExecName + " <command> [<args>]")
    fmt.Println()
    fmt.Println("Available commands include:")
    fmt.Println("   decode        Converts a binary data/event log file to a readable CSV file")
    fmt.Println("   help          Displays information about using this utility")
    fmt.Println("   info          Displays the information about this utility and what's installed")
    fmt.Println("   install       Install or upgrade the Strongback Java Library")
    fmt.Println("   new-project   Creates a new project configured to use Strongback (only 1.x)")
    fmt.Println("   releases      Display the available versions of the Strongback Java Library")
    fmt.Println("   version       Display the currently installed version")
    fmt.Println("   uninstall     Remove an existing Strongback Java Library installation")
    fmt.Println()
    fmt.Println("Additional help is available for each command with:")
    fmt.Println()
    fmt.Println("   " + ExecName + " help <command>")
    fmt.Println()
}

func PrintInstallUsage() {
    fmt.Println("   " + ExecName + " install [--skip-archive] [--verbose] [version] ")
    fmt.Println()
    fmt.Println("Description:")
    fmt.Println("   Install or upgrade the Strongback Java Library.")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println()
    fmt.Println("   --skip-archive")
    fmt.Println("       Do not archive the current installation before installing the new version.")
    fmt.Println("       This flag does nothing if there is no current installation.")
    fmt.Println()
    fmt.Println("   --verbose")
    fmt.Println("       Print additional detailed information during the operation.")
    fmt.Println()
    fmt.Println("Arguments:")
    fmt.Println()
    fmt.Println("   version")
    fmt.Println("       The version of the Strongback Java Library that should be installed.")
    fmt.Println("       The latest version is used if an explicit version is not provided.")
    fmt.Println()
}

func PrintUninstallUsage() {
    fmt.Println("   " + ExecName + " uninstall [--skip-archive] [--remove-archives] [--verbose] [--yes]")
    fmt.Println()
    fmt.Println("Description:")
    fmt.Println("   Remove any Strongback Java Library that is already installed.")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println()
    fmt.Println("   --skip-archive")
    fmt.Println("       Do not archive the current installation before removing it.")
    fmt.Println("       This flag does nothing if there is no current installation.")
    fmt.Println()
    fmt.Println("   --remove-archives")
    fmt.Println("       Remove all archives of previous Strongback Java Library installations.")
    fmt.Println("       Doing this is permanent and will prevent recoverying previous installations.")
    fmt.Println()
    fmt.Println("   --verbose")
    fmt.Println("       Print additional detailed information during the operation.")
    fmt.Println()
    fmt.Println("   --yes")
    fmt.Println("       Do not prompt about removing existing installation or archives.")
    fmt.Println()
}

func PrintReleasesUsage() {
    fmt.Println("   " + ExecName + " releases [--all]")
    fmt.Println()
    fmt.Println("Description:")
    fmt.Println("   List the releases Strongback Java Library that are available as listed on the")
    fmt.Println("   " + ExecName + " GitHub organization. The installed version is highlighted.")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println("   --all")
    fmt.Println("       Show all the releases, including alpha, beta, and other early")
    fmt.Println("       releases that are only for testing purposes. Unless this is")
    fmt.Println("       provided, only releases that are ready for use on robots are listed.")
    fmt.Println()
}

func PrintNewProjectUsage() {
    fmt.Println("   " + ExecName + " new-project [--directory <path>] [--package <packageName>]")
    fmt.Println("                          [--no-eclipse] [--overrite]")
    fmt.Println("                          name")
    fmt.Println()
    fmt.Println("Description:")
    fmt.Println("     Create a new FRC robot project using the Strongback Java Library.")
    fmt.Println("     No files will be overwritten unless --overwrite is specified.")
    fmt.Println()
    fmt.Println("Arguments:")
    fmt.Println()
    fmt.Println("   name")
    fmt.Println("       The name of the new project.")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println("   --dir <parent_directory>")
    fmt.Println("       The directory where this utility should place the new project.")
    fmt.Println("       Defaults to the current directory.")
    fmt.Println()
    fmt.Println("   --package")
    fmt.Println("       Specifies a custom initial package for Robot.java. Defaults to 'org.frc<teamNumber>.robot'")
    fmt.Println("       where the team number is obtained from the WPILib installation or is '0' if WPILib is not")
    fmt.Println("       installed and initialized through Eclipse.")
    fmt.Println()
    fmt.Println("   --no-eclipse")
    fmt.Println("       Use this if you are not using Eclipse to avoid creating Eclipse project metadata files.")
    fmt.Println()
    fmt.Println("   --overwrite")
    fmt.Println("       Forces overwriting of existing files.")
    fmt.Println()
}

func PrintDecodeUsage() {
    fmt.Println("   " + ExecName + " decode [--verbose] input [output]")
    fmt.Println()
    fmt.Println("Description:")
    fmt.Println("   Converts binary log files to readable CSV files")
    fmt.Println()
    fmt.Println("Arguments:")
    fmt.Println()
    fmt.Println("   input")
    fmt.Println("       The path to the binary log recorded on the robot by the Strongback Java Library.")
    fmt.Println()
    fmt.Println("   output")
    fmt.Println("       The path to the file to be written by this utility and that will contain the")
    fmt.Println("       comma separated values (CSV). If not provided, the output will be saved in the")
    fmt.Println("       current directory in a file with the same filename as the input but with a .csv extension.")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println()
    fmt.Println("   --verbose")
    fmt.Println("       Print additional detailed information during the operation.")
    fmt.Println()
}

func PrintVersionUsage() {
    fmt.Println("   " + ExecName + " version")
    fmt.Println()
    fmt.Println("Description:")
    fmt.Println("   Output the shortened version information for the Strongback utility and Java Library")
    fmt.Println()
}

func PrintInfoUsage() {
    fmt.Println("   " + ExecName + " info")
    fmt.Println()
    fmt.Println("Description:")
    fmt.Println("   Output the detailed version information for the Strongback utility and Java Library")
    fmt.Println()
}

func PrintUsageError(err error) {
    fmt.Println()
    fmt.Println("Error: " + err.Error())
}

func PrintUsageLead() {
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Println()
}

func HasFlagsAfterArguments(command *flag.FlagSet) bool {
    for i := 0; i!=command.NArg(); i++ {
        arg := command.Arg(i)
        if strings.HasPrefix(arg, "--") {
            fmt.Println()
            fmt.Println("Error: unexpected flag " + arg + " appearing after arguments")
            return true          
        }
    }
    return false
}

func main() {
    // Subcommands without flags
    versionCommand := flag.NewFlagSet("version", flag.ContinueOnError)
    helpCommand := flag.NewFlagSet("help", flag.ContinueOnError)

    // Subcommands with flags
    releasesCommand := flag.NewFlagSet("releases", flag.ContinueOnError)
    allReleases := releasesCommand.Bool("all", false, "List all releases.")

    installCommand := flag.NewFlagSet("install", flag.ContinueOnError)
    installSkipArchive := installCommand.Bool("skip-archive", false, "Do not create an archive before upgrading.")
    installVerbose := installCommand.Bool("verbose", false, "Print additional detail.")

    uninstallCommand := flag.NewFlagSet("uninstall", flag.ContinueOnError)
    uninstallSkipArchive := uninstallCommand.Bool("skip-archive", false, "Do not create an archive before removing.")
    uninstallVerbose := uninstallCommand.Bool("verbose", false, "Print additional detail.")
    uninstallYes := uninstallCommand.Bool("yes", false, "Do not prompt to remove.")
    uninstallRemoveArchives := uninstallCommand.Bool("remove-archives", false, "Also remove all archives.")

    decodeCommand := flag.NewFlagSet("decode", flag.ContinueOnError)
    decodeVerbose := decodeCommand.Bool("verbose", false, "Print additional detail.")

    newProjectCommand := flag.NewFlagSet("new-project", flag.ContinueOnError)
    newProjectNoEclipse := newProjectCommand.Bool("no-eclipse", false, "Avoid generating Eclipse metadata for project.")
    newProjectOverwrite := newProjectCommand.Bool("overwrite", false, "Overwrite existing files.")
    newProjectDirectory := newProjectCommand.String("directory", "", "Directory.")
    newProjectPackage := newProjectCommand.String("package", "", "Directory.")

    // Verify that a subcommand has been provided
    // os.Arg[0] is the main command
    // os.Arg[1] will be the subcommand
    if len(os.Args) < 2 {
        PrintUsage()
        os.Exit(1)
    }

    // Switch on the subcommand
    // Parse the flags for appropriate FlagSet
    // FlagSet.Parse() requires a set of arguments to parse as input
    // os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
    switch os.Args[1] {
    case "install":
        installCommand.SetOutput(bytes.NewBuffer([]byte{}))
        if err := installCommand.Parse(os.Args[2:]); err != nil {
            PrintUsageError(err)
            PrintUsageLead()
            PrintInstallUsage()
            os.Exit(1)
        }
        if HasFlagsAfterArguments(installCommand) {
            PrintUsageLead()
            PrintInstallUsage()
            os.Exit(1)
        }
        var desiredVersion string
        if installCommand.NArg() > 0 {
            desiredVersion = installCommand.Arg(0)
        }
        env := NewEnvironment()
        env.InstallRelease(desiredVersion, *installSkipArchive, *installVerbose)
        os.Exit(0)
    case "uninstall":
        uninstallCommand.SetOutput(bytes.NewBuffer([]byte{}))
        if err := uninstallCommand.Parse(os.Args[2:]); err != nil {
            PrintUsageError(err)
            PrintUsageLead()
            PrintUninstallUsage()
            os.Exit(1)
        }
        if HasFlagsAfterArguments(uninstallCommand) {
            PrintUsageLead()
            PrintInstallUsage()
            os.Exit(1)
        }
        env := NewEnvironment()
        env.UninstallRelease(*uninstallYes, *uninstallSkipArchive, *uninstallVerbose, *uninstallRemoveArchives)
        os.Exit(0)
    case "decode":
        decodeCommand.SetOutput(bytes.NewBuffer([]byte{}))
        if err := decodeCommand.Parse(os.Args[2:]); err != nil {
            PrintUsageError(err)
            PrintUsageLead()
            PrintDecodeUsage()
            os.Exit(1)
        }
        if HasFlagsAfterArguments(decodeCommand) || decodeCommand.NArg() == 0 {
            PrintUsageLead()
            PrintInstallUsage()
            os.Exit(1)
        }
        inputFile := decodeCommand.Arg(0)
        if !files.IsExistingFile(inputFile) {
            fmt.Println("Unable to locate binary log file: " + inputFile)
            os.Exit(4)
        }
        outputFile := ""
        if decodeCommand.NArg() > 1 {
            outputFile = decodeCommand.Arg(1)
        }
        env := NewEnvironment()
        env.CheckInstalled()
        env.DecodeFile(inputFile, outputFile, *decodeVerbose)
        os.Exit(0)
    case "new-project":
        newProjectCommand.SetOutput(bytes.NewBuffer([]byte{}))
        if err := newProjectCommand.Parse(os.Args[2:]); err != nil {
            PrintUsageError(err)
            PrintUsageLead()
            PrintNewProjectUsage()
            os.Exit(1)
        }
        if HasFlagsAfterArguments(newProjectCommand) || newProjectCommand.NArg() == 0 {
            PrintUsageLead()
            PrintNewProjectUsage()
            os.Exit(1)
        }
        projectName := newProjectCommand.Arg(0)
        directory, err := filepath.Abs("./")
        if err != nil {
            log.Fatal(err)
            os.Exit(1)
        }
        if len(*newProjectDirectory) > 0 {
            directory = *newProjectDirectory
        }
        javaPackage := ""
        if len(*newProjectPackage) > 0 {
            javaPackage = *newProjectPackage            
        }
        if files.IsExistingFileOrDirectory(directory + files.PathSeparator + projectName) {
            fmt.Println()
            fmt.Println("Error: directory already exists and --overwrite not provided")
            PrintUsageLead()
            PrintNewProjectUsage()
            os.Exit(3)
        }
        env := NewEnvironment()
        env.CheckInstalled()
        if env.NewProject(projectName, directory, javaPackage, !*newProjectNoEclipse, *newProjectOverwrite, false) {
            os.Exit(0)
        }
        os.Exit(4)
    case "releases":
        releasesCommand.SetOutput(bytes.NewBuffer([]byte{}))
        if err := releasesCommand.Parse(os.Args[2:]); err != nil {
            PrintUsageError(err)
            PrintUsageLead()
            PrintReleasesUsage()
            os.Exit(1)
        }
        env := NewEnvironment()
        env.PrintReleases(*allReleases)
        os.Exit(0)
    case "version":
        versionCommand.SetOutput(bytes.NewBuffer([]byte{}))
        if err := versionCommand.Parse(os.Args[2:]); err != nil {
            PrintUsageError(err)
            PrintUsageLead()
            PrintVersionUsage()
            os.Exit(1)
        }
        env := NewEnvironment()
        env.PrintVersion()
        os.Exit(0)
    case "info":
        env := NewEnvironment()
        env.PrintInfo()
        os.Exit(0)
    case "help":
        helpCommand.SetOutput(bytes.NewBuffer([]byte{}))
        helpCommand.Parse(os.Args[2:])
        PrintUsageLead()
        var command string
        if helpCommand.NArg() > 0 {
            command = helpCommand.Arg(0)
            switch command {
            case "install":
                PrintInstallUsage()
            case "uninstall":
                PrintUninstallUsage()
            case "decode":
                PrintDecodeUsage()
            case "new-project":
                PrintNewProjectUsage()
            case "releases":
                PrintReleasesUsage()
            case "version":
                PrintVersionUsage()
            case "info":
                PrintInfoUsage()
            default:
                PrintUsage()
            }
        } else {
            PrintUsage()
        }
        os.Exit(0)
    default:
        PrintUsageLead()
        PrintUsage()
        os.Exit(1)
    }

	// dir := UserHomeDir()
    env := NewEnvironment()
    env.Print()
}


