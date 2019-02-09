Feature: Built-ins

Background:
   Given I'm in project dir

Scenario: Run other tasks
  Given file ./Oyafile containing
    """
    Project: project

    baz: |
      echo "baz"

    bar: |
      echo "bar"
      $Tasks.baz()
    """
  When I run "oya run bar"
  Then the command succeeds
  And the command outputs to stdout
  """
  bar
  baz

  """

Scenario: Run pack's tasks
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    bar: |
      echo "bar"
      $Tasks.baz()

    baz: |
      echo "baz"
    """
  When I run "oya run foo.bar"
  Then the command succeeds
  And the command outputs to stdout
  """
  bar
  baz

  """

Scenario: Pack can only run its own tasks
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1
      github.com/test/bar: v0.0.1

    Import:
      foo: github.com/test/foo
      bar: github.com/test/bar
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    foo: |
      echo "foo"
    """
  And file ./.oya/packs/github.com/test/bar@v0.0.1/Oyafile containing
    """
    bar: |
      $Tasks.foo()
    """
  When I run "oya run bar.bar"
  Then the command fails with error matching
    """"
    .*variable not found.*
    """"

Scenario: Access Oyafile base directory
  Given file ./Oyafile containing
    """
    Project: project
    """
  And file ./subdir/Oyafile containing
    """
    all: |
      echo $BasePath
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  ^.*subdir

  """

Scenario: Access pack base directory
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    all: |
      echo $BasePath
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  ^.*github.com/test/foo@v0.0.1

  """

Scenario: Access Oyafile Project name
  Given file ./Oyafile containing
    """
    Project: project

    all: |
      echo $Project
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  project

  """

Scenario: Access Oyafile Project name in nested dir
  Given file ./Oyafile containing
    """
    Project: project
    """
  And file ./subdir/Oyafile containing
    """
    all: |
      echo $Project
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  project

  """

Scenario: Access Oyafile Project name inside pack
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    all: |
      echo $Project
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  project

  """
