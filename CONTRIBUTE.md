The [Strongback](http://strongback.org) project is operated as a community-centric open source project. Although founded by [FRC Team 4931](http://evilletech.com/), we hope that many teams want to use Strongback and contribute to our community. Everyone is welcome in our community. See also the [Strongback Community page](https://github.com/strongback/strongback-java/wiki/Community).

**The Strongback CLI and Strongback Java Library are separate projects.** The [Strongback Java Library](https://github.com/strongback/strongback-java) is updated fairly frequently and is used on FRC robots; see https://github.com/strongback/strongback-java for more details. 

The rest of this document describes how to contribute to the **Strongback CLI**, which can install and use multiple versions of the Strongback Java Library.

### Fork the Strongback CLI repository

Go to the [Strongback CLI repository](https://github.com/strongback/strongback-cli) and press the "Fork" button near the upper right corner of the page. When finished, you will have your own "fork" at `https://github.com/<your-github-username>/strongback-cli`, and this is the repository to which you will upload your proposed changes and create pull requests. For details, see the [GitHub documentation](https://help.github.com/articles/fork-a-repo/).

### Building Locally

The Strongback CLI is written in [Go](http://golang.org), so first make sure you've [installed](https://golang.org/doc/install) the Go development tools, set up a [workspace](https://golang.org/doc/code.html#Workspaces), and [set the GOPATH](https://golang.org/doc/install#testing) environment variable.

Next, make the `src/strongback.org` directories inside your workspace, and open a terminal in that `strongback.org` directory. Run the following commmand to clone the Strongback CLI repository using HTTPS authentication::

    $ git clone https://github.com/<your-github-username>/strongback-cli.git cli

or if you prefer to use SSH and have [uploaded your public key to your GitHub account](https://help.github.com/articles/adding-a-new-ssh-key-to-your-github-account/), you can instead use SSH:

    $ git clone git@github.com:<your-github-username>/strongback-cli.git cli

This will create a `cli` directory, so change into that directory:

	$ cd cli

This local repository knows about your fork, but it doesn't yet know about the official or "upstream" Strongback CLI repository. So, run the following commands:

	$ git remote add upstream https://github.com/strongback/strongback-cli.git
	$ git fetch upstream
	$ git branch --set-upstream-to=upstream/master master

Now, when you check the status using Git, it will compare your local repository to the official _upstream_ repository.

### Get the latest upstream code

You will frequently need to get all the of the changes that are made to the upstream repository, and you can do this with these commands:

    $ git fetch upstream
    $ git pull upstream master

The first command fetches all changes on all branches, while the second actually updates your local `master` branch with the latest commits from the `upstream` repository.

### Building locally

To build the source code locally, checkout and update the `master` branch:

    $ git checkout master
    $ git pull upstream master

Then compile and package the CLI code for all of the platforms:

    $ make all

This command will compile and build all executables and place them in OS-specific directories inside the `out` directory (e.g., `out/macos/strongback`, `out/linux/strongback`, and `out/windows/strongback.exe`). It will also create the `.tar.gz` and `.zip` archives for each of the platforms. 

As you make changes, you probably want to only compile for your own platform:

    $ make clean out/macos

or

    $ make clean out/linux


### Running and debugging tests

To run the compiled CLI tool, simply run the `strongback` (or `strongback.exe`) executable for your platform. For example, run the `help` command on OS using:

    $ out/macos/strongback help

or on Linux with:

    $ out/linux/strongback help

### Making changes

Everything the community does with the codebase -- fixing bugs, adding features, making improvements, adding tests, etc. -- should be described by an issue in our [issue tracker](https://github.com/strongback/strongback-cli/issues). If no such issue exists for what you want to do, please create an issue with a meaningful and easy-to-understand description.

Before you make any changes, be sure to switch to the `master` branch and pull the latest commits on the `master` branch from the upstream repository. Also, it's probably good to run a build and verify all tests pass *before* you make any changes.

    $ git checkout master
    $ git pull upstream master
    $ mvn clean install

Once everything builds, create a *topic branch* named appropriately (we recommend using the issue number, such as `issue-1234`):

    $ git checkout -b issue-1234

This branch exists locally and it is there you should make all of your proposed changes for the issue. As you'll soon see, it will ultimately correspond to a single pull request that the Strongback committers will review and merge (or reject) as a whole. (Some issues are big enough that you may want to make several separate but incremental sets of changes. In that case, you can create subsequent topic branches for the same issue by appending a short suffix to the branch name.)

Please verify your changes compile and work before committing them. Feel free to commit your changes locally as often as you'd like, though we generally prefer that each commit represent a complete and atomic change to the code. Often, this means that most issues will be addressed with a single commit in a single pull-request, but other more complex issues might be better served with a few commits that each make separate but atomic changes. (Some developers prefer to commit frequently and to ammend their first commit with additional changes. Other developers like to make multiple commits and to then squash them. How you do this is up to you. However, *never* change, squash, or ammend a commit that appears in the history of the upstream repository.) When in doubt, use a few separate atomic commits; if the Strongback reviewers think they should be squashed, they'll let you know when they review your pull request.

Committing is as simple as:

    $ git commit .

which should then pop up an editor of your choice in which you should place a good commit message. _*We do expect that all commit messages' first line ends starts with a short phrase that summarizes what changed in the commit and ends with the issue number in parentheses.*_ For example:

    Corrected help documentation (Issue #1234)

If that phrase is not sufficient to explain your changes, then the first line should be followed by a blank line and one or more paragraphs with additional details. 

### Rebasing

If its been more than a day or so since you created your topic branch, we recommend *rebasing* your topic branch on the latest `master` branch. This requires switching to the `master` branch, pulling the latest changes, switching back to your topic branch, and rebasing:

    $ git checkout master
    $ git pull upstream master
    $ git checkout issue-1234
    $ git rebase master

If your changes are compatible with the latest changes on `master`, this will complete and there's nothing else to do. However, if your changes affect the same files/lines as other changes have since been merged into the `master` branch, then your changes conflict with the other recent changes on `master`, and you will have to resolve them. The git output will actually tell you you need to do (e.g., fix a particular file, stage the file, and then run `git rebase --continue`), but if you have questions consult Git or GitHub documentation or spend some time reading about Git rebase conflicts on the Internet.

### Creating a pull request

Once you're finished making your changes, your topic branch should have your commit(s) and you should have verified that your branch builds and runs successfully. At this point, you can shared your proposed changes and create a pull request. To do this, first push your topic branch (and its commits) to your fork repository (called `origin`) on GitHub:

    $ git push origin issue-1234

Then, in a browser go to https://github.com/strongback/strongback-cli, and you should see a small section near the top of the page with a button labeled "Create pull request". GitHub recognized that you pushed a new topic branch to your fork of the upstream repository, and it knows you probably want to create a pull request with those changes. Click on the button, and GitHub will present you with a short form that you should fill out with information about your pull request. The title should start with the issue number and include a short summary of the changes included in the pull request. (If the pull request contains a single commit, GitHub will automatically prepopulate the title and description fields from the commit message.) Add a description with details about your change and end the description with:

    Fixes #1234

where `1234` is the issue number to which this pull request corresponds. When completed, press the "Create" button and copy the URL to the new pull request.

At this point, you can switch to another issue and another topic branch. The Strongback committers will be notified of your new pull request, and will review it in short order. They may ask questions or make remarks using line notes or comments on the pull request. (By default, GitHub will send you an email notification of such changes, although you can control this via your GitHub preferences.)

If the reviewers ask you to make additional changes, simply switch to your topic branch for that pull request:

    $ git checkout issue-1234

and then make the changes on that branch and either add a new commit or ammend your previous commits. When you've addressed the reviewers' concerns, push your changes to your `origin` repository:

    $ git push origin issue-1234

GitHub will automatically update the pull request with your latest changes, but we ask that you go to the pull request and add a comment summarizing what you did. This process may continue until the reviewers are satisfied.

By the way, please don't take offense if the reviewers ask you to make additional changes, even if you think those changes are minor. The reviewers have a broach understanding of the codebase, and their job is to ensure the code remains as uniform as possible and of sufficient quality. When they believe your pull request has those attributes, they will merge your pull request into the official upstream repository.

Once your pull request has been merged, feel free to delete your topic branch both in your local repository:

    $ git branch -d issue-1234

and in your fork: 

    $ git push origin :issue-1234

(This last command is a bit strange, but it basically is pushing an empty branch (the space before the `:` character) to the named branch. Pushing an empty branch is the same thing as removing it.)


### Summary

Here's a quick check list for a good pull request:

* An existing issue that describes the problem, enhancement, or new feature
* One commit per PR, with the issue number in the commit comment
* One feature/change per PR
* No changes to code not directly related to your change (e.g. no formatting changes or refactoring to existing code, if you want to refactor/improve existing code that's a separate discussion and separate issue)
* A full build completes succesfully
* Do a rebase on upstream `master`