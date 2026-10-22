# Graph Report - app-listener  (2026-10-22)

## Corpus Check
- 519 files · ~1,736,677 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6768 nodes · 29941 edges · 245 communities (197 shown, 48 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2766 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `286b5ae1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- globreserve_test.go
- go_pkg_path_filepath
- go_pkg_strings
- protectedPage
- find_test.go
- github.com/charmbracelet/bubbletea.Cmd
- installPage
- IntegrationSuite
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- BinaryEntry
- guard.bpf.c
- sync.Mutex
- Finder
- ListUsers
- hashCache
- recoverFileInPlace
- TestMonitorUseCaseLifecycle
- gui.go
- codeedit.go
- picker
- CandidateDir
- IntegrationSuite
- Message
- controlServer
- TempBinary
- Network
- IntegrationSuite
- absPath
- monitorUnitTest
- daemon_multicall_test.go
- fileedit_model_test.go
- expandPlaceholders
- os.File
- fadeScroll
- steamrt_test.go
- cursor.go
- Client
- daemonconfig.go
- navItem
- Plan
- planUpdaters
- mutedText
- engine
- go_pkg_os
- New
- Recorder
- exe_supersede.h
- updatePage
- recordFlow
- Config
- context.Context
- IntegrationSuite
- .FilterExistingWhitelist
- liveSession
- DaemonUseCase
- newTestTarget
- Guard
- fcntl
- bpfstats/main.go
- FileTarget
- codecheck/main.go
- rcUser
- githubRelease
- process_vm_readv.c
- fakeGuardRepo
- LibraryClosure
- SelectOrphans
- runtimeKeeper
- net_tester/main.go
- github.com/charmbracelet/bubbles/textarea.Model
- Binary
- io.Writer
- dialLiveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- EventType
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- ConfinedEntry
- pinstate.go
- Well-known bypass classes
- DumpHook
- jit_provenance.c
- inspector_probe.c
- recordops.c
- runGuard
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- GuardInodeKey
- mountTable
- gui/session.go
- NetGuardEvent
- IntegrationSuite
- github.com/spf13/cobra.Command
- secureResources
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- GenerateConf
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- pins.go
- buildTrustedSet
- trace-app-libs.sh
- IntegrationSuite
- vaultfs.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- runDaemon
- serve.go
- dialog_test.go
- appTheme
- Ledger
- editauth_test.go
- go_pkg_github_com_testcontainers_testcontainers_go
- Vault
- ValidateTargets
- fakeStore
- inodeChain
- IntegrationSuite
- buildOneGuard
- safeio.go
- lockAndDeprovision
- catalogRefresher
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- IntegrationSuite
- .populateWatchInodes
- feed.go
- shell
- editView
- time.Time
- WithBottomBar
- findUpdate
- extractor
- Limitations
- .install
- IntegrationSuite
- Active
- netinject.c
- copyRegularAt
- inode_dev.h
- run
- Compatibility
- How it works
- searchDirs
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- testKeyPair
- RawTarget
- edit-protected
- mustSubkey
- kernel-prof.sh
- update_test.go
- pathCache
- systemd_test.go
- DaemonEvent
- changelog_test.go
- Daemon configuration
- uninstall/uninstall.go
- Desktop GUI
- IntegrationSuite
- procstats.go
- IntegrationSuite
- .eventTypeAllowed
- bleedLayout
- releasesSince
- NewDumpHook
- .awaitDaemonUpOrDead
- TestCleanRel
- openBunTmp
- tempGrantSetup
- configedit_test.go
- IntegrationSuite
- ptrace_race.c
- prctl_gate
- compileGlobNames
- audit.go
- supersede_probe.c
- bpfLsmListed
- IntegrationSuite
- fakeDaemon
- Vault
- github.com/Virgula0/app-listener/tools/codecheck
- wm.go
- .runEditorHarness
- taint_victim.c
- parsePasswd
- Troubleshooting
- logLine
- parseGuardEvent
- wm_test.go
- AddServeFlags
- logGrid
- github.com/stretchr/testify/suite.Suite
- Execute
- monitor: see what touches a path

## God Nodes (most connected - your core abstractions)
1. `Guard` - 122 edges
2. `Config` - 116 edges
3. `GuardInodeKey` - 102 edges
4. `shell` - 91 edges
5. `Load()` - 90 edges
6. `absPath()` - 88 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go

## Import Cycles
- None detected.

## Communities (245 total, 48 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.16
Nodes (28): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+20 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (131): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+123 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.19
Nodes (12): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+4 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (45): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_user(), fill_identity() (+37 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (22): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+14 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.15
Nodes (37): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+29 more)

### Community 10 - "go_pkg_path_filepath"
Cohesion: 0.35
Nodes (4): recordTick, isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 11 - "go_pkg_strings"
Cohesion: 0.06
Nodes (22): releaseBuild(), Run(), atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape() (+14 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (23): backendLink, checkGrid, configView, forwardForm, forwardView, protectedPage, recordForm, recordReview (+15 more)

### Community 13 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 14 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (42): recordTickCmd(), recordModel, diffModel, editorModel, formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision() (+34 more)

### Community 15 - "installPage"
Cohesion: 0.07
Nodes (14): actionGate, binaryRow, binSource, counterStrip, handCheck, iconTile, installAction, installPage (+6 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.10
Nodes (56): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+48 more)

### Community 18 - "runView"
Cohesion: 0.04
Nodes (25): daemonPage, InstallMode, InstallRequest, installUI, Options, runHost, runView, UninstallHooks (+17 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (8): backingDeviceUnion(), Resource, BackingDevices(), GuardRepository, encryptionRootOf(), uniqueEncryptionRoots(), vaultOpForGuard(), vaultOpOnRoot()

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "BinaryEntry"
Cohesion: 0.08
Nodes (20): deferredBinary, keptEntries, refusedReplacements, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys(), canonicalBinaryPath(), commMatchesGuardedBinary() (+12 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (73): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+65 more)

### Community 26 - "sync.Mutex"
Cohesion: 0.05
Nodes (29): forgetHooks, tempGrant, trustRows, fyneLogFilter, freedSuperseded(), PruneSuperseded(), setSupersedeTrust(), supersedeMaps() (+21 more)

### Community 27 - "Finder"
Cohesion: 0.15
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 28 - "ListUsers"
Cohesion: 0.11
Nodes (28): runFacultatives(), cleanOrphanedFscrypt(), dupFD(), installSSHAgent(), installSSHAgentEnv(), TestInstallSSHAgentRefusesSymlinkedUnitDir(), removeFacultative(), removeFacultatives() (+20 more)

### Community 29 - "hashCache"
Cohesion: 0.07
Nodes (29): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, placedWalk (+21 more)

### Community 30 - "recoverFileInPlace"
Cohesion: 0.20
Nodes (18): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+10 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "gui.go"
Cohesion: 0.17
Nodes (16): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), stagedBinary, Mode (+8 more)

### Community 33 - "codeedit.go"
Cohesion: 0.10
Nodes (43): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+35 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (21): keyEntry, listOpts, pickEntry, picker, treeSearch, lstatOf(), outsideHardlinks(), TestSearchTreeFollowsNoSymlink() (+13 more)

### Community 35 - "CandidateDir"
Cohesion: 0.13
Nodes (24): selectAndEditConfig(), addManualDirectories(), groupCandidates(), pickDirectories(), pickFromCandidates(), pickUsers(), TestGroupCandidates(), Candidate (+16 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (12): checkSecret(), TestLiveSessionWatchAndPinger(), holdForward(), pingInterval(), File, backend, backendOp, backendState (+4 more)

### Community 38 - "controlServer"
Cohesion: 0.11
Nodes (21): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+13 more)

### Community 39 - "TempBinary"
Cohesion: 0.12
Nodes (16): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+8 more)

### Community 40 - "Network"
Cohesion: 0.04
Nodes (28): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), netGuardPill(), BinariesSummary(), TestBinariesSummary() (+20 more)

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (8): IntegrationSuite, IntegrationSuite, IntegrationSuite, IntegrationSuite, absPath(), TestMain(), helperChild(), TestMain()

### Community 44 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "expandPlaceholders"
Cohesion: 0.09
Nodes (16): reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), expandPlaceholders(), BinaryRule (+8 more)

### Community 48 - "os.File"
Cohesion: 0.04
Nodes (75): multicallAdmissible(), checkFreshDir(), openDirNoFollow(), renewDir(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), setConfinedHomes() (+67 more)

### Community 49 - "fadeScroll"
Cohesion: 0.07
Nodes (20): Boot ordering, Lifecycle, Reload (`systemctl reload`), Shutdown and crashes, Startup order, fadeBar, fadeBarRenderer, fadeScroll (+12 more)

### Community 50 - "steamrt_test.go"
Cohesion: 0.14
Nodes (21): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), openDir(), TestApplyPinsReplacesOpaqueAndMergesOthers(), TestPinsStale() (+13 more)

### Community 51 - "cursor.go"
Cohesion: 0.08
Nodes (22): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+14 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (13): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Hello (+5 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.12
Nodes (41): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+33 more)

### Community 54 - "navItem"
Cohesion: 0.04
Nodes (13): clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabButton, tabButtonRenderer (+5 more)

### Community 55 - "Plan"
Cohesion: 0.12
Nodes (15): Plan, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), TestPlanAdmitsOnlyWhatWasRecorded(), TestPlanExtendsAnExistingSection() (+7 more)

### Community 56 - "planUpdaters"
Cohesion: 0.15
Nodes (17): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+9 more)

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (39): binaryList, filesPage, homePage, multiAnswer, panel, newLogo(), newBinaryList(), newBinaryRow() (+31 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (28): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), TestSetInspectorsBounded(), nearestRootClaims(), pathWithin() (+20 more)

### Community 59 - "go_pkg_os"
Cohesion: 0.07
Nodes (48): asRoot(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256() (+40 more)

### Community 60 - "New"
Cohesion: 0.24
Nodes (16): looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip(), TestVaultUnlockLockDispatchForRegularFile() (+8 more)

### Community 61 - "Recorder"
Cohesion: 0.17
Nodes (14): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+6 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.16
Nodes (18): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+10 more)

### Community 63 - "updatePage"
Cohesion: 0.23
Nodes (7): installedBinary, latestRelease, updatePage, readInstalled(), newUpdatePage(), boolIndex(), isPreRelease()

### Community 64 - "recordFlow"
Cohesion: 0.14
Nodes (17): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+9 more)

### Community 65 - "Config"
Cohesion: 0.08
Nodes (52): entryWriters(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), addedResources() (+44 more)

### Community 66 - "context.Context"
Cohesion: 0.18
Nodes (10): reachable(), centered(), FitSize(), poll(), TestFitSize(), hyprClient, hyprMonitor, Rect (+2 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 69 - "liveSession"
Cohesion: 0.22
Nodes (10): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck(), normalizedEvents() (+2 more)

### Community 70 - "DaemonUseCase"
Cohesion: 0.16
Nodes (3): ConcurrencyLimit(), DaemonUseCase, unlockRoots()

### Community 71 - "newTestTarget"
Cohesion: 0.10
Nodes (22): NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst() (+14 more)

### Community 72 - "Guard"
Cohesion: 0.05
Nodes (27): holdHeadless(), HeadlessLine(), lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), rootHandle, snapGuard, openRoot() (+19 more)

### Community 73 - "fcntl"
Cohesion: 0.09
Nodes (11): main(), report(), denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "FileTarget"
Cohesion: 0.05
Nodes (32): netGuardTab, policyEditor, policyHooks, policyTab, specFromForm(), guardPill(), newPolicyTab(), verbFor() (+24 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "rcUser"
Cohesion: 0.12
Nodes (40): rcTarget, BunTmpDir(), TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv() (+32 more)

### Community 78 - "githubRelease"
Cohesion: 0.21
Nodes (18): changelogText(), cumulativeChangelog(), applyUpdate(), archAssets(), assetsFor(), deploy(), downloadAndVerify(), fetchGUI() (+10 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.19
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 81 - "LibraryClosure"
Cohesion: 0.11
Nodes (19): hostView(), isScript(), LibraryClosure(), mountVouchesOwnership(), rootOwnedInode(), rootOwnedSafe(), TestLibraryClosureRejectsUserWritableOriginDir(), runsNatively() (+11 more)

### Community 82 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 83 - "runtimeKeeper"
Cohesion: 0.27
Nodes (8): ownerPaths(), keptRuntimes(), ownedBy(), startRuntimeKeeper(), TestKeptRuntimesNeedsTheUsersReadOnlyLibDir(), keptRuntime, runtimeKeeper, OnBeforeForgetFreed()

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (19): lineStyles, pos, viewportState, reposition(), fg(), baseStyle(), gutterLabel(), lineLen() (+11 more)

### Community 86 - "Binary"
Cohesion: 0.10
Nodes (27): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), classifyLibs(), finishBinary() (+19 more)

### Community 87 - "io.Writer"
Cohesion: 0.20
Nodes (11): formatElapsed(), orNone(), printBinary(), printRecording(), candidates(), confirmable(), confirmAll(), describe() (+3 more)

### Community 88 - "dialLiveSession"
Cohesion: 0.06
Nodes (36): auditAfterEdit(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply(), editAgain() (+28 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "EventType"
Cohesion: 0.07
Nodes (16): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+8 more)

### Community 92 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.04
Nodes (45): lsmHook, trustHook, foreignPinnedLink(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+37 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.14
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 100 - "ConfinedEntry"
Cohesion: 0.12
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 101 - "pinstate.go"
Cohesion: 0.08
Nodes (41): configureDaemonLogging(), loadDaemonConfig(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder() (+33 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "DumpHook"
Cohesion: 0.24
Nodes (5): logHook, formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (14): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+6 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "runGuard"
Cohesion: 0.07
Nodes (32): DrainUntilSignal(), CheckBPFLSM(), CheckEBPF(), confirmMasterKeyOverwrite(), newPinGeneration(), prepareDaemonStart(), runBPFCheck(), runGenKey() (+24 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 110 - "GuardInodeKey"
Cohesion: 0.05
Nodes (31): exeRow, inodeBatchFn, tempMaskOp, tempRow, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches() (+23 more)

### Community 111 - "mountTable"
Cohesion: 0.20
Nodes (8): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), ReadMountinfo(), unescapeMountinfo(), mountTable

### Community 112 - "gui/session.go"
Cohesion: 0.18
Nodes (16): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+8 more)

### Community 113 - "NetGuardEvent"
Cohesion: 0.09
Nodes (19): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode() (+11 more)

### Community 115 - "github.com/spf13/cobra.Command"
Cohesion: 0.20
Nodes (15): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials(), TestSetupDumpLogRefusesSymlink() (+7 more)

### Community 116 - "secureResources"
Cohesion: 0.24
Nodes (8): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.04
Nodes (77): logRefreshChange(), TestLogRefreshChange(), collectFilesystemPrereqs(), askEncryption(), refreshCatalogConfig(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse() (+69 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "pins.go"
Cohesion: 0.20
Nodes (11): systemTool(), isWhiteout(), mountStage(), opaque(), pinRecord(), readPinDir(), readPins(), writeNew() (+3 more)

### Community 131 - "buildTrustedSet"
Cohesion: 0.14
Nodes (17): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+9 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "vaultfs.go"
Cohesion: 0.12
Nodes (23): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead() (+15 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "runDaemon"
Cohesion: 0.07
Nodes (27): newBinaryVetter(), openBinaryVetter(), newControlManager(), startControlManager(), makeReloadHandler(), markReady(), noteReloadedSettings(), notifySystemdReady() (+19 more)

### Community 143 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 144 - "dialog_test.go"
Cohesion: 0.39
Nodes (11): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+3 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "editauth_test.go"
Cohesion: 0.11
Nodes (25): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+17 more)

### Community 148 - "go_pkg_github_com_testcontainers_testcontainers_go"
Cohesion: 0.15
Nodes (4): attrOpCase, inNosuidNamespace(), exploitTest, guardExploitTest

### Community 149 - "Vault"
Cohesion: 0.23
Nodes (3): Vault, Open(), TestOpenVaultRefusesSymlink()

### Community 150 - "ValidateTargets"
Cohesion: 0.16
Nodes (12): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets(), mustParse() (+4 more)

### Community 152 - "inodeChain"
Cohesion: 0.23
Nodes (15): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+7 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.09
Nodes (38): buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon(), startGuardedDaemonAbortable(), systemPatterns() (+30 more)

### Community 155 - "safeio.go"
Cohesion: 0.25
Nodes (13): invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), removeCheckedFile(), AtomicWriteAt(), CreateExclNoFollow(), descend() (+5 more)

### Community 156 - "lockAndDeprovision"
Cohesion: 0.13
Nodes (11): deprovisionKind, migrationTarget, writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr(), classifyMigrationTarget(), isDeprovisionBusy(), isDeprovisionMissing() (+3 more)

### Community 157 - "catalogRefresher"
Cohesion: 0.12
Nodes (14): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps() (+6 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 165 - "feed.go"
Cohesion: 0.07
Nodes (41): actionBar, cell, column, dataTable, eventFeed, feedAction, feedSpec, networkPage (+33 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (48): pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines, snapNet (+40 more)

### Community 167 - "editView"
Cohesion: 0.11
Nodes (11): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+3 more)

### Community 168 - "time.Time"
Cohesion: 0.10
Nodes (17): LockedError, addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+9 more)

### Community 169 - "WithBottomBar"
Cohesion: 0.31
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 170 - "findUpdate"
Cohesion: 0.21
Nodes (14): installHooks(), staged(), channelOf(), FetchLatest(), FetchVerified(), filterChannel(), findUpdate(), Release (+6 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - ".install"
Cohesion: 0.10
Nodes (21): Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Process inspectors and taint, Quiet metadata denials, Related, Self-protection, Steam runtime updates (+13 more)

### Community 175 - "Active"
Cohesion: 0.16
Nodes (12): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, EditText() (+4 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 178 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 179 - "run"
Cohesion: 0.18
Nodes (14): checkUpdatePreconditions(), deployGUI(), InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), RunWith() (+6 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 186 - "Binary trust"
Cohesion: 0.25
Nodes (8): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, The dynamic loader, trust

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 195 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+2 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 199 - "update_test.go"
Cohesion: 0.10
Nodes (17): downloadFile(), downloadReleaseFiles(), fetchRelease(), fetchReleases(), newerThanInstalled(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchGUIWithoutAsset() (+9 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "systemd_test.go"
Cohesion: 0.05
Nodes (57): anyDaemonRunning(), offerDaemonStart(), applyLiveRefresh(), buildBinaryIfNeeded(), installBinaryOnly(), markInstallComplete(), preflightDeployedBinary(), runRefreshCatalogOnly() (+49 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.08
Nodes (38): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+30 more)

### Community 203 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 204 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 205 - "uninstall/uninstall.go"
Cohesion: 0.10
Nodes (35): deletePostBackups(), listBackups(), restoreBackups(), askDecrypt(), confirmKeyDeletion(), decryptDirectories(), decryptStep(), offerBackupCleanup() (+27 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.20
Nodes (10): Daemon, Desktop GUI, File and folder pickers, Files, Home, Network, Protected, Tables (+2 more)

### Community 209 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "releasesSince"
Cohesion: 0.36
Nodes (9): newReleaseWindow(), releasesSince(), compareStableVersions(), latestTime(), newerThanStable(), parseStableVersion(), parseVerPart(), releaseWindow (+1 more)

### Community 214 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 217 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 218 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 219 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 221 - "ptrace_race.c"
Cohesion: 0.52
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 223 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 224 - "audit.go"
Cohesion: 0.08
Nodes (40): auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+32 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 226 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 231 - "wm.go"
Cohesion: 0.36
Nodes (6): Detect(), Session, hyprlandSocket(), i3RuntimeSocket(), socketOwnedBy(), Kind

### Community 234 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "logLine"
Cohesion: 0.25
Nodes (5): logLine, journalLine(), TestJournalLine(), logRow(), logStyle()

### Community 238 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 239 - "wm_test.go"
Cohesion: 0.67
Nodes (8): envOf(), serve(), shortDir(), TestDetect(), TestHyprlandFloat(), TestHyprlandFloatTimesOutWithoutWindow(), TestI3Float(), uid()

### Community 240 - "AddServeFlags"
Cohesion: 0.43
Nodes (7): AddHeadlessFlag(), AddServeFlags(), init(), init(), init(), init(), init()

### Community 243 - "Execute"
Cohesion: 0.67
Nodes (3): Execute(), initLogger(), main()

### Community 244 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 618 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **48 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `go_pkg_strings`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `go_pkg_github_com_testcontainers_testcontainers_go`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.020540089667396517 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `sync.Mutex`, `engine`, `fakeDaemon`, `time.Duration`, `ConfinedEntry`, `guardUnitTest`, `time.Time`, `DaemonEvent`, `runGuard`, `runDaemon`, `GuardInodeKey`, `os.File`, `BinaryEntry`, `buildOneGuard`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09124209124209125 - nodes in this community are weakly interconnected._
- **Why does `daemon: permanent protection with encryption at rest` connect `.install` to `fadeScroll`, `README.md`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07737874861162532 - nodes in this community are weakly interconnected._