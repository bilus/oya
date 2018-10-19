# Oya

## Usage

1. Install oya and its dependencies:

        curl https://raw.githubusercontent/bilus/oya/master/scripts/setup.sh | sh

1. Initialize project to use a certain CI/CD tool and workflow. Example:

        oya init jenkins-monorepo

   It boostraps configuration for Jenkins pipelines supporting the 1.a workflow (see Workflows below), an Oyafile and supporting scripts and compatible generators.

1. Run a hook:

        oya run build

   Right now it won't do anything as there are no buildable directories yet. Let's create one.

1. Create a buildable directory:

        mkdir app1
        cat > Oyafile
        build: echo "Hello, world"

1. Run the build hook again:

        oya run build
        Hello, world


## How it works

A directory is included in the build process if it has an Oyafile. Let's call that such directory a
"buildable directory".

Oya first walks all directories to build the changeset: a list of buildable directories.
It then walks the list of directories, running the requested hook for each directory

Hooks and their corresponding scripts are defined in `Oyafile`s. Names of hooks can be arbitrary camel-case yaml identifiers, starting with a lower-case letter. Built-in hooks start with capital letters.

Example `Oyafile`:

```
build: docker build .
test: pytest
```

## Changesets

   * `changeset` -- (optional) modifies the current changeset (see Changesets).

Oya first walks all directories to build the changeset: a list of buildable directories.
It then walks the list, running the matching hook in each.
   CI/CD tool-specific script outputting list of modified files in buildable directories given the current hook name.
     - each path must be normalized and prefixed with `+`
     - cannot be overriden, only valid for top-level Oyafile
     - in the future, you'll be able to override for a buildable directory and use `-` to exclude directories, `+` to include additional ones,
       and use wildcards, this will allow e.g. forcing running tests for all apps when you change a shared directory
     - git diff --name-only origin/master...$branch
     - https://dzone.com/articles/build-test-and-deploy-apps-independently-from-a-mo
     - https://stackoverflow.com/questions/6260383/how-to-get-list-of-changed-files-since-last-build-in-jenkins-hudson/9473207#9473207

Generation of the changeset is controlled by the optional changeset key in Oyafiles,
which can point to a script executed to generate the changeset:

1. No directive -- includes all directories containing on Oyafile.
2. Directive pointing to a script.

.oyaignore lists files whose changes do not trigger build for the containing buildable directory

## Features/ideas

1. Generators based on packs. https://github.com/Flaque/thaum + draft pack plugin

## Workflows

### Repo structure

1. Mono-repo:
   - Each app has its own directory
   - There is a directory/file containing deployment configuration

2. Multi-repo:
   - Each app has its own repo
   - Deployment configurations in its own repo

3. Mix:
   - Some/all apps share repos, some may have their own
   - Deployment configurations in its own repo

> Also submodules tried for NT/Switchboard and eventually ditched.

### Change control

a. Each environment has its own directory

b. Each environment has its own branch

### Evaluation

| Workflow | Projects    | Pros                                       | Cons                                                    |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 1.a      | E           | "Can share code"                           | Merge order dependent [1]                               |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 1.b      | C           | Single checkout                            | Complex automation [2]                                  |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 2.a      |             | Same as 2.b                                | Same as 2.b plus need to detect which directory changed |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 2.b      | S           | Better isolation [3] Simple automation [4] | More process overhead [5]                               |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 3.a      | P           | Can divide up a project however you like   | Complex automation [2]                                  |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 3.b      | P           | Simple deployment automation               | Same as 3.a                                             |

* [1] Code gets merged from branch to branch; works for small team.
* [2] Need to detect what changed between commits. Many CI/CD tools allow only one configuration per repo and require coding around the limitations, example: https://discuss.circleci.com/t/does-circleci-2-0-work-with-monorepos/10378/13
* [3] No way to just share code, need to package into libraries. Great for microservices and must have for large teams.
* [4] Just put a CI/CD config into the root.
* [5] No way to just share code, need to package into libraries. Bad for small teams wanting to quickly prototype.
