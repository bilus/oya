Feature: Built-ins

Background:
   Given I'm in project dir

# Tests for invoking other tasks via $Tasks, do not actually invoke
# oya but rather output the command line. See oya_test.go (MustSetUp)
# for details.

@xxx
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
  oya run baz

  """

@xxx
Scenario: Run pack's tasks
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
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
  oya run foo.baz

  """

# TODO: Use tasks in packs.
# TODO: BasePath.
