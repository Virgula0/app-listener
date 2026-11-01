# Graph Report - app-listener  (2026-11-01)

## Corpus Check
- 565 files · ~1,797,775 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7379 nodes · 33767 edges · 251 communities (208 shown, 43 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3011 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c18865c7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- op
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- GuardEvent
- options.go
- binaryList
- OpenConfined
- buildOneGuard
- recorder/plan.go
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- FileTarget
- guard.bpf.c
- LibraryClosure
- BinaryEntry
- io.Writer
- .showSized
- codeedit.go
- TestMonitorUseCaseLifecycle
- fileedit_model_test.go
- bufio.Reader
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- Selection
- watchSet
- daemon.conf reference
- showConfirm
- monitorUnitTest
- Candidate
- fileedit_model.go
- Guard
- stdio
- journal
- fadeScroll
- IntegrationSuite
- SanitizeText
- Client
- installPage
- sync.Mutex
- .apply
- backend_test.go
- mutedText
- engine
- go_pkg_os
- runNetworkGuard
- Serve
- root.go
- runOfflineEdit
- Network
- unix_connect.c
- wm_test.go
- IntegrationSuite
- Recorder
- FindDaemonProcesses
- Vault
- LaunchRule
- planUpdaters
- fcntl
- New
- update_test.go
- codecheck/main.go
- recordFlow
- newTestVetter
- process_vm_readv.c
- os.File
- Binary
- .filterEvent
- rcUser
- golang.org/x/sys/unix.Stat_t
- bpfstats/main.go
- exe_supersede.h
- Config
- vaultfs_test.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- tempGrantSetup
- maintainerKeeper
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- SampleContent
- NetworkMonitor
- Well-known bypass classes
- guardModel
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- CandidateDir
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- .posAt
- updatePage
- UI
- promptPassword
- IntegrationSuite
- ConfinedEntry
- systemd.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- newTestTarget
- IntegrationSuite
- WebSocket /ws client
- inet_connect.c
- install.sh
- hashCache
- classify
- make test-integration (rootful Docker suite)
- backend.go
- Gate
- SecretStoreNotes
- trace-app-libs.sh
- IntegrationSuite
- Monitor
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- serveWebSocket
- liveSession
- IntegrationSuite
- Ledger
- FileEvent
- TempBinary
- RawTarget
- Daemon configuration
- fakeStore
- mc_tag.h
- IntegrationSuite
- github.com/charmbracelet/bubbles/textarea.Model
- catalogRefresher
- CopyTreeWithProgress
- EventType
- Installation
- IntegrationSuite
- graphify-refresh.sh
- .Layout
- testKeyPair
- dialog_test.go
- absPath
- feed.go
- shell
- inodeChain
- plan_test.go
- Vault
- Classify
- NetworkMonitorRepository
- Limitations
- appTheme
- parseDaemonEvents
- TrustGuard
- net_tester/main.go
- NetEventType
- githubRelease
- sources
- .netIntegrityContainer
- fileedit_symlink_test.go
- runView
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- harnessSuite
- os.FileInfo
- changelog_test.go
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- controlManager
- pathCache
- pinstate_test.go
- time.Time
- sync.Once
- GuardInodeKey
- AddServeFlags
- monitor
- IntegrationSuite
- run
- IntegrationSuite
- launchKey
- KernelDev
- WithBottomBar
- linkView
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- Multicall
- github.com/charmbracelet/bubbletea.Cmd
- app-listener documentation
- Plan
- restore.go
- entryOf
- VouchedDevs
- editView
- supersede_probe.c
- Run
- bleedLayout
- mountTable
- go_pkg_strings
- github.com/Virgula0/app-listener/tools/codecheck
- inode_dev.h
- network monitor: see what a program talks to
- loadDaemonConfig
- ValidateTargets
- TestMain
- Troubleshooting
- deploy
- GenerateConf
- sources_test.go
- IntegrationSuite
- taint_victim.c
- User
- .runEditorHarness
- parsePasswd
- bootstrapSelfProtectPlaceholders
- prctl_gate
- monitor: see what touches a path
- TestCleanRel
- daemonconfig.go

## God Nodes (most connected - your core abstractions)
1. `Config` - 147 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `shell` - 97 edges
5. `Load()` - 97 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 75 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Modes` --references--> `monitor()`  [INFERRED]
  docs/README.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (251 total, 43 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.19
Nodes (31): tempJournalRow, argvMatch, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion (+23 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (161): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+153 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote(), watchPathsInConfig()

### Community 3 - "op"
Cohesion: 0.16
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (83): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+75 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (18): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+10 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+27 more)

### Community 10 - "GuardEvent"
Cohesion: 0.04
Nodes (20): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+12 more)

### Community 11 - "options.go"
Cohesion: 0.36
Nodes (9): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+1 more)

### Community 12 - "binaryList"
Cohesion: 0.10
Nodes (14): binaryList, netGuardTab, policyEditor, policyHooks, policyTab, newBinaryList(), newPolicyTab(), verbFor() (+6 more)

### Community 13 - "OpenConfined"
Cohesion: 0.05
Nodes (43): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+35 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.06
Nodes (44): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), displayEntries() (+36 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.08
Nodes (65): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+57 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (23): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+15 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (51): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+43 more)

### Community 22 - "Resource"
Cohesion: 0.04
Nodes (18): multicallAdmissible(), ownerPaths(), scriptIdentity(), backingDeviceUnion(), findResource(), fakeDaemon, vetDecision, Resource (+10 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "FileTarget"
Cohesion: 0.15
Nodes (9): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+1 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (95): inode_dev(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+87 more)

### Community 26 - "LibraryClosure"
Cohesion: 0.08
Nodes (27): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+19 more)

### Community 27 - "BinaryEntry"
Cohesion: 0.10
Nodes (20): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath() (+12 more)

### Community 28 - "io.Writer"
Cohesion: 0.22
Nodes (6): formatElapsed(), orNone(), printBinary(), printRecording(), fyneLogFilter, limitedWriter

### Community 29 - ".showSized"
Cohesion: 0.13
Nodes (13): installUI, multiAnswer, sizedDialog, scrollable(), newFadeScroll(), ask(), diffText(), orDefault() (+5 more)

### Community 30 - "codeedit.go"
Cohesion: 0.11
Nodes (40): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), moveTo(), Navigate() (+32 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.13
Nodes (14): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+6 more)

### Community 32 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 33 - "bufio.Reader"
Cohesion: 0.11
Nodes (19): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+11 more)

### Community 34 - "picker"
Cohesion: 0.09
Nodes (19): keyEntry, listOpts, pickEntry, picker, treeSearch, TestSearchTreeFollowsNoSymlink(), errText(), isFindShortcut() (+11 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.22
Nodes (25): busPath(), connexec(), connFD(), createItem(), dial(), execInto(), failf(), fakebus() (+17 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (14): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+6 more)

### Community 38 - "controlServer"
Cohesion: 0.15
Nodes (9): authSelfPeer(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+1 more)

### Community 39 - "Selection"
Cohesion: 0.28
Nodes (4): pos, cursorPos(), Selection, offset()

### Community 40 - "watchSet"
Cohesion: 0.18
Nodes (9): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, `[secret <name>]`, Sections at a glance (+1 more)

### Community 41 - "daemon.conf reference"
Cohesion: 0.09
Nodes (22): A complete example, `allow_env`, `allow_lib`, Binary lines and events, daemon.conf reference, `deny_args`, `warn_args`, `deny_env`, `warn_env`, `[electron_apps]` (+14 more)

### Community 42 - "showConfirm"
Cohesion: 0.12
Nodes (12): buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+4 more)

### Community 43 - "monitorUnitTest"
Cohesion: 0.13
Nodes (4): sameInode(), dirTarget(), newMountTable(), monitorUnitTest

### Community 44 - "Candidate"
Cohesion: 0.09
Nodes (33): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), appendSectionsAndEdit(), collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), addManualDirectories() (+25 more)

### Community 45 - "fileedit_model.go"
Cohesion: 0.11
Nodes (18): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), groupWatchPaths() (+10 more)

### Community 46 - "Guard"
Cohesion: 0.05
Nodes (18): deferredBinary, exeRow, rootHandle, supersededBinary, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), SharedPinDegraded(), TestNoteRootMountFollowsMountID() (+10 more)

### Community 48 - "journal"
Cohesion: 0.13
Nodes (7): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (30): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+22 more)

### Community 50 - "IntegrationSuite"
Cohesion: 0.07
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 51 - "SanitizeText"
Cohesion: 0.05
Nodes (25): backendLink, configView, forwardForm, forwardView, handRadioGroup, protectedPage, recordView, trustPage (+17 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (14): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+6 more)

### Community 53 - "installPage"
Cohesion: 0.11
Nodes (9): actionGate, binSource, installAction, installPage, uninstallCard, optionHint(), warnLabel(), daemonStopNote() (+1 more)

### Community 54 - "sync.Mutex"
Cohesion: 0.10
Nodes (9): peerCred(), TestReadFrame(), decodeFrame(), encodeWithSerial(), header(), readFrame(), queued, senderQueue (+1 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (15): Apply(), DeleteBackup(), Outcome, ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+7 more)

### Community 56 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (44): checkGrid, filesPage, handCheck, homePage, recordForm, recordReview, newLogo(), newBinaryRow() (+36 more)

### Community 58 - "engine"
Cohesion: 0.05
Nodes (26): engine, portOwner, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), pathWithin(), TestPathWithinMatchesPrefixRule() (+18 more)

### Community 59 - "go_pkg_os"
Cohesion: 0.07
Nodes (37): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin() (+29 more)

### Community 60 - "runNetworkGuard"
Cohesion: 0.08
Nodes (33): DrainUntilSignal(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runMonitor() (+25 more)

### Community 61 - "Serve"
Cohesion: 0.16
Nodes (17): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+9 more)

### Community 62 - "root.go"
Cohesion: 0.06
Nodes (33): TestParseReservedPorts(), releaseBuild(), Run(), closeDumpFile(), Execute(), initLogger(), postRun(), Progress (+25 more)

### Community 63 - "runOfflineEdit"
Cohesion: 0.13
Nodes (15): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), TestPickEncryptedDirSingleSkipsPicker(), chooseResource() (+7 more)

### Community 64 - "Network"
Cohesion: 0.11
Nodes (8): netGuardPill(), NetworkGuardRepository, stopPumped(), clonePolicy(), NetExtra, NetPolicy, Network, verifyHashes()

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm_test.go"
Cohesion: 0.11
Nodes (24): reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket(), poll() (+16 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.21
Nodes (6): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 68 - "Recorder"
Cohesion: 0.21
Nodes (12): Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), TestMonitorTargetsNameFilesByAbsolutePath(), TestProgressNamesMulticallApplets() (+4 more)

### Community 69 - "FindDaemonProcesses"
Cohesion: 0.10
Nodes (20): applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), refreshCatalogConfig(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr(), TestUnifiedDiff(), TestUnifiedDiffIdentical() (+12 more)

### Community 70 - "Vault"
Cohesion: 0.06
Nodes (43): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), TestDiffMergePreservesExistingAndParses(), encryptDirectories(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs() (+35 more)

### Community 71 - "LaunchRule"
Cohesion: 0.15
Nodes (17): binaryLaunchRules(), launchRulesByPath(), TestLaunchRulesByPath(), argvEntry, applyLaunchRule(), envName(), LaunchRule, ParseAllowEnv() (+9 more)

### Community 72 - "planUpdaters"
Cohesion: 0.10
Nodes (26): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), hardLinked(), libBinaries(), mayUpdate() (+18 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "New"
Cohesion: 0.08
Nodes (49): decryptDirectories(), TestDecryptDirectoriesEmpty(), classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord() (+41 more)

### Community 75 - "update_test.go"
Cohesion: 0.09
Nodes (19): downloadFile(), downloadReleaseFiles(), fetchRelease(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), TestFetchReleases() (+11 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "recordFlow"
Cohesion: 0.15
Nodes (12): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+4 more)

### Community 78 - "newTestVetter"
Cohesion: 0.15
Nodes (19): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+11 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.05
Nodes (64): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+56 more)

### Community 81 - "Binary"
Cohesion: 0.13
Nodes (25): lines(), pickBinaries(), pickLibs(), reviewRecording(), appletLine(), beyondEvents(), classifyGroup(), classifyLibs() (+17 more)

### Community 83 - "rcUser"
Cohesion: 0.12
Nodes (40): rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir() (+32 more)

### Community 84 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.11
Nodes (23): allowedEnv(), TestAllowedEnvKeepsOnlyRootPlaced(), placedWalk, treeJudge, BuiltinAllowedEnv(), confineBelow(), inUserHome(), GlobSystemPlaced() (+15 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.08
Nodes (38): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+30 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.20
Nodes (16): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+8 more)

### Community 87 - "Config"
Cohesion: 0.04
Nodes (72): newBinaryVetter(), openBinaryVetter(), newCatalogRefresher(), buildGuards(), makeReloadHandler(), markReady(), markReadyWithLogin(), notifySystemdReady() (+64 more)

### Community 88 - "vaultfs_test.go"
Cohesion: 0.16
Nodes (17): DefaultNewFileMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead(), mustWrite(), openRootT() (+9 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 92 - "maintainerKeeper"
Cohesion: 0.19
Nodes (10): appRunning(), libDirSig(), grantNewBinaries(), maintainedApp, maintainerKeeper, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending() (+2 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.06
Nodes (24): foreignPinnedLink(), WithSuspectsFrom(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), codeVouched(), isTraced() (+16 more)

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
Nodes (14): procStats, statsSample, duration(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "SampleContent"
Cohesion: 0.18
Nodes (15): TestDaemonUnitWithMetadataOutput(), installFileAs(), installUnit(), upsertFile(), isInstallerSSHAgentUnit(), TestIsInstallerSSHAgentUnit(), TestDaemonUnitCmdPrintsBundledUnit(), SampleContent() (+7 more)

### Community 101 - "NetworkMonitor"
Cohesion: 0.09
Nodes (10): snapNetMonitor, NetEvent, BinariesSummary(), TestBinariesSummary(), PumpRingbuf(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure() (+2 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "guardModel"
Cohesion: 0.06
Nodes (36): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+28 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "CandidateDir"
Cohesion: 0.05
Nodes (43): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+35 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - ".posAt"
Cohesion: 0.09
Nodes (16): lineStyles, viewportState, fg(), TestDaemonConfTokens(), TestTokenizeKeepsLineCount(), tokenClassOf(), baseStyle(), gutterLabel() (+8 more)

### Community 111 - "updatePage"
Cohesion: 0.09
Nodes (17): anyDaemonRunning(), handCheckGroup, handSelect, installedBinary, latestRelease, updatePage, handCursor(), newCheckGroup() (+9 more)

### Community 112 - "UI"
Cohesion: 0.08
Nodes (25): nextConfig(), installConfig(), showChangelog(), channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith() (+17 more)

### Community 113 - "promptPassword"
Cohesion: 0.20
Nodes (12): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), promptPassword(), promptPasswordWith() (+4 more)

### Community 115 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 116 - "systemd.go"
Cohesion: 0.09
Nodes (43): offerDaemonStart(), buildBinaryIfNeeded(), certifyInstall(), ensureInstalledBinary(), installBinaryOnly(), markInstallComplete(), preflightDeployedBinary(), prepareInstallation() (+35 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 127 - "classify"
Cohesion: 0.22
Nodes (12): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify(), classifyCollection() (+4 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "backend.go"
Cohesion: 0.12
Nodes (20): askDecrypt(), confirmKeyDeletion(), decryptStep(), pickDirsToDecrypt(), run(), runUninstall(), RunWith(), TestRunWithRefusesConcurrentRun() (+12 more)

### Community 130 - "Gate"
Cohesion: 0.16
Nodes (7): placedKey(), Caller, Gate, Role, Listen(), Bus, Judge

### Community 131 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "Monitor"
Cohesion: 0.12
Nodes (7): NewMonitor(), putBtrfsLayout(), TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace(), Monitor, watchedInode

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 143 - "serveWebSocket"
Cohesion: 0.15
Nodes (12): configureWebSocket(), decodeEndsAtEOF(), parseResize(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestParseResize(), TestSameOrigin() (+4 more)

### Community 144 - "liveSession"
Cohesion: 0.15
Nodes (16): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward() (+8 more)

### Community 145 - "IntegrationSuite"
Cohesion: 0.23
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (27): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+19 more)

### Community 147 - "FileEvent"
Cohesion: 0.12
Nodes (6): snapMonitor, FileEvent, MonitorRepository, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 148 - "TempBinary"
Cohesion: 0.10
Nodes (20): confirmUserPlaced(), forwardCheck(), TempBinary, TemporaryGrant, TempRule, planTempBlock(), ResolveTempBinary(), TestResolveTempBinaryClassesRuntime() (+12 more)

### Community 149 - "RawTarget"
Cohesion: 0.24
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 150 - "Daemon configuration"
Cohesion: 0.15
Nodes (12): `[inspectors]`, Daemon configuration, Daemon settings, Generic Electron apps, Launch rules, Library trust, Port sections, Process inspectors (+4 more)

### Community 152 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 154 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.09
Nodes (31): findButton, span, editorModel, column(), findAll(), Finder, lowerRunes(), NewFinder() (+23 more)

### Community 155 - "catalogRefresher"
Cohesion: 0.06
Nodes (31): drain(), logRefreshChange(), readEvents(), sameConfig(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps() (+23 more)

### Community 156 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "EventType"
Cohesion: 0.08
Nodes (19): selfEventsFor(), secretJudge, SecretScope, BpfEvent, EventType, eventUnitTest, parseEvents(), SecretMatch (+11 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 163 - "dialog_test.go"
Cohesion: 0.39
Nodes (11): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+3 more)

### Community 164 - "absPath"
Cohesion: 0.06
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (47): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+39 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (57): fakeLink, pageDef, pageID, reel, resizeWatch, shell, snapEngines, snapNet (+49 more)

### Community 167 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 168 - "plan_test.go"
Cohesion: 0.27
Nodes (14): isWhiteout(), opaque(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles() (+6 more)

### Community 169 - "Vault"
Cohesion: 0.21
Nodes (4): Vault, Open(), TestOpenVaultRefusesSymlink(), Unopened()

### Community 170 - "Classify"
Cohesion: 0.18
Nodes (7): Classify(), drain(), fixtureExe(), record(), TestClassifyVerifiesTheInodeThatRan(), TestObserveKeysByKernelExeAndSkipsSelf(), harnessSuite

### Community 171 - "NetworkMonitorRepository"
Cohesion: 0.18
Nodes (4): NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase(), StreamUseCase

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 174 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 175 - "TrustGuard"
Cohesion: 0.05
Nodes (27): argvMatches, lsmHook, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, deleteResKeys() (+19 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "NetEventType"
Cohesion: 0.25
Nodes (8): NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType()

### Community 178 - "githubRelease"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 179 - "sources"
Cohesion: 0.16
Nodes (8): TestGeneratedLibraryBlockRoundTrips(), lstatOf(), outsideHardlinks(), pathByte(), pathRecord(), TestPathRecord(), libDirDev(), sources

### Community 181 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 182 - "runView"
Cohesion: 0.05
Nodes (14): daemonPage, logGrid, logHook, logLine, runHost, runView, journalLine(), TestJournalLine() (+6 more)

### Community 186 - "Binary trust"
Cohesion: 0.22
Nodes (9): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch rules, Launch scan, Multicall binaries, The dynamic loader (+1 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "github.com/spf13/cobra.Command"
Cohesion: 0.14
Nodes (23): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink(), applyVerbosity() (+15 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 194 - "os.FileInfo"
Cohesion: 0.19
Nodes (6): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 195 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "controlManager"
Cohesion: 0.33
Nodes (3): newControlManager(), startControlManager(), controlManager

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "pinstate_test.go"
Cohesion: 0.07
Nodes (42): CheckBPFLSM(), CheckEBPF(), configureDaemonLogging(), confirmMasterKeyOverwrite(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck() (+34 more)

### Community 202 - "time.Time"
Cohesion: 0.07
Nodes (41): LockedError, ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI() (+33 more)

### Community 204 - "GuardInodeKey"
Cohesion: 0.05
Nodes (41): inodeBatchFn, refusedReplacements, TestNearestRootClaims(), TestTieRank(), TestSetInspectorsBounded(), appendInodeRow(), batchInodeRows(), inodeRowsOf() (+33 more)

### Community 205 - "AddServeFlags"
Cohesion: 0.36
Nodes (8): AddHeadlessFlag(), AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init()

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (33): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+25 more)

### Community 209 - "run"
Cohesion: 0.36
Nodes (8): InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestReadInstalledVersion(), GUIInstalled()

### Community 212 - "KernelDev"
Cohesion: 0.09
Nodes (29): BusClient, fdStat, BinaryStat, physicalParent(), statFD(), AttestBusClient(), busClientKey(), execStamp() (+21 more)

### Community 213 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 214 - "linkView"
Cohesion: 0.36
Nodes (4): linkTarget(), entryKind, libPlan, linkView

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "Multicall"
Cohesion: 0.29
Nodes (7): Multicall, MulticallKind, multicallFamily(), classifyBytes(), fwd(), TestClassifyMulticall(), TestMulticallMarkersReversed()

### Community 218 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (19): recordTickCmd(), recordModel, snapNetGuard, diffModel, newDiffModel(), NetGuardEvent, listenForGuardEvents(), tickGuardStats() (+11 more)

### Community 219 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 220 - "Plan"
Cohesion: 0.47
Nodes (4): dedupe(), Plan, libOfChosen(), Extension

### Community 221 - "restore.go"
Cohesion: 0.09
Nodes (31): runDiffCatalog(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups() (+23 more)

### Community 222 - "entryOf"
Cohesion: 0.40
Nodes (5): entryOf(), FileModeToUnixMode(), TestFileModeToUnixModeRoundTrip(), TestUnixModeToFileMode(), UnixModeToFileMode()

### Community 223 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 224 - "editView"
Cohesion: 0.11
Nodes (9): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+1 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "Run"
Cohesion: 0.33
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+2 more)

### Community 229 - "go_pkg_strings"
Cohesion: 0.09
Nodes (5): parseReservedPorts(), readReservedPorts(), resizeMessage, i3Result, i3Workspace

### Community 231 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 232 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 233 - "loadDaemonConfig"
Cohesion: 0.22
Nodes (9): loadDaemonConfig(), noteReloadedSettings(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds(), SetUserHomes() (+1 more)

### Community 234 - "ValidateTargets"
Cohesion: 0.18
Nodes (10): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets(), TestValidateTargets() (+2 more)

### Community 235 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "deploy"
Cohesion: 0.13
Nodes (27): installHooks(), staged(), archAssets(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), deploy(), downloadAndVerify() (+19 more)

### Community 238 - "GenerateConf"
Cohesion: 0.06
Nodes (60): LibraryBlock, RefreshOptions, SectionScan, IsScript(), TestIsScript(), SecretServiceBlock(), AppendSectionWhitelist(), EnsureSecretServiceBlock() (+52 more)

### Community 239 - "sources_test.go"
Cohesion: 0.38
Nodes (4): newSources(), TestSourcesHoldFilesAndArchiveMembers(), TestSourcesRefuseOtherFilesystem(), writeTar()

### Community 242 - "User"
Cohesion: 0.06
Nodes (67): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), addedResources(), applyDiffAdditions() (+59 more)

### Community 244 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 245 - "bootstrapSelfProtectPlaceholders"
Cohesion: 1.00
Nodes (3): ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders()

### Community 246 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 247 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 267 - "daemonconfig.go"
Cohesion: 0.08
Nodes (65): systemPatterns(), systemRule(), watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBlockLine(), applyConfigLine() (+57 more)

## Knowledge Gaps
- **33 isolated node(s):** `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 680 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **43 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `go_pkg_strings`, `TestMain`, `parseDaemonEvents`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.netIntegrityContainer`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01868552052521463 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `app-listener documentation`, `Installation`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0943859649122807 - nodes in this community are weakly interconnected._
- **Why does `Installation` connect `Installation` to `monitor`, `README.md`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06569133677567413 - nodes in this community are weakly interconnected._