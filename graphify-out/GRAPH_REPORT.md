# Graph Report - app-listener  (2026-10-28)

## Corpus Check
- 526 files · ~1,744,616 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6830 nodes · 30248 edges · 232 communities (191 shown, 41 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2774 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `729a24f1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- wipe
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- globreserve_test.go
- fakeGuardRepo
- New
- protectedPage
- os.File
- go_pkg_path_filepath
- installPage
- IntegrationSuite
- guard_trust.bpf.c
- tabs.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- Guard
- guard.bpf.c
- sync.Mutex
- Finder
- FileTarget
- .step
- github.com/charmbracelet/bubbletea.Cmd
- TestMonitorUseCaseLifecycle
- binaryList
- find_test.go
- picker
- fscrypt/orphans.go
- IntegrationSuite
- Message
- controlServer
- codeedit.go
- daemon_multicall_test.go
- IntegrationSuite
- absPath
- monitorUnitTest
- highlight.go
- fileedit_model_test.go
- journal
- stdio
- runOfflineEdit
- fadeScroll
- vaultfs.go
- cursor.go
- proto.go
- daemonconfig.go
- runView
- .apply
- mutedText
- GuardInodeKey
- Vault
- os.FileInfo
- Recorder
- exe_supersede.h
- compileGlobNames
- FileEvent
- github.com/charmbracelet/bubbles/textarea.Model
- context.Context
- IntegrationSuite
- classifyGroup
- liveSession
- editView
- batchInodeRows
- Candidate
- string
- bpfstats/main.go
- lockAndDeprovision
- codecheck/main.go
- shellrc.go
- update_test.go
- process_vm_readv.c
- GuardEvent
- op
- User
- maintainerKeeper
- net_tester/main.go
- updatePage
- Binary
- planUpdaters
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
- NetGuardEvent
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- .populateWatchInodes
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- hashFileShared
- mountTable
- wm.go
- recordFlow
- IntegrationSuite
- procstats.go
- Open
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
- install/install.go
- LibraryClosure
- trace-app-libs.sh
- IntegrationSuite
- time.Time
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- catalogRefresher
- Serve
- NewDumpHook
- UI
- Ledger
- configedit_test.go
- Daemon configuration
- changelog_test.go
- Config
- fakeStore
- IntegrationSuite
- IntegrationSuite
- buildOneGuard
- testKeyPair
- fakeDaemon
- fileedit_symlink_test.go
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- RawTarget
- IntegrationSuite
- parseGuardEvent
- feed.go
- shell
- newTestTarget
- watchSet
- SanitizeText
- wm_test.go
- runLockdown
- Limitations
- runBackend
- IntegrationSuite
- github.com/spf13/cobra.Command
- netinject.c
- CopyTreeWithProgress
- stripTempJournal
- TestMain
- Compatibility
- How it works
- KernelDev
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- writeWithin
- IntegrationSuite
- edit-protected
- plan_test.go
- kernel-prof.sh
- MonitorRepository
- pathCache
- bpfLsmListed
- DaemonEvent
- Plan
- loadDaemonConfig
- RunWith
- Desktop GUI
- IntegrationSuite
- guardModel
- Prereq
- .eventTypeAllowed
- bleedLayout
- inode_dev.h
- prctl_gate
- .awaitDaemonUpOrDead
- Vault
- watchPattern
- tempGrantSetup
- harnessSuite
- IntegrationSuite
- BinaryEntry
- FileEngines
- runWizard
- supersede_probe.c
- inspector_probe.c
- links_test.go
- github.com/Virgula0/app-listener/tools/codecheck
- taint_victim.c
- Active
- Troubleshooting
- monitor: see what touches a path

## God Nodes (most connected - your core abstractions)
1. `Guard` - 123 edges
2. `Config` - 118 edges
3. `GuardInodeKey` - 103 edges
4. `shell` - 91 edges
5. `Load()` - 90 edges
6. `absPath()` - 89 edges
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

## Communities (232 total, 41 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.10
Nodes (35): Execute(), initLogger(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs (+27 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (142): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity() (+134 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "wipe"
Cohesion: 0.19
Nodes (10): isFileVaultCiphertext(), isRegularFileTarget(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall(), TestNewBoundedKeyFnRetryAborts() (+2 more)

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
Cohesion: 0.05
Nodes (21): addErrKind, TestWithChmodDropWriteNeedsReadOnly(), canonicalPaths(), classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), modeString(), NewGuard() (+13 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.11
Nodes (45): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify() (+37 more)

### Community 10 - "fakeGuardRepo"
Cohesion: 0.07
Nodes (16): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant(), newExeTap() (+8 more)

### Community 11 - "New"
Cohesion: 0.09
Nodes (43): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+35 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (29): backendLink, checkGrid, configView, forwardForm, forwardView, installUI, multiAnswer, protectedPage (+21 more)

### Community 13 - "os.File"
Cohesion: 0.04
Nodes (76): checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes(), ReplacementCheck, VettedInode, Multicall, MulticallKind (+68 more)

### Community 14 - "go_pkg_path_filepath"
Cohesion: 0.31
Nodes (6): nextConfig(), recordTick, isAbsoluteCandidate(), TestCatalogSanity(), EditText(), hint

### Community 15 - "installPage"
Cohesion: 0.06
Nodes (16): actionGate, binaryRow, binSource, clickPanel, handButton, handCheck, iconTile, installAction (+8 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (60): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+52 more)

### Community 18 - "tabs.go"
Cohesion: 0.03
Nodes (18): appTheme, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar, tabButton (+10 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (10): clipLabel(), humanSize(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode, fileEditModel (+2 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.16
Nodes (46): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+38 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (11): Resource, GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), uniqueEncryptionRoots() (+3 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (38): CheckPassword(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+30 more)

### Community 24 - "Guard"
Cohesion: 0.05
Nodes (27): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, deferredBinary, exeRow, rootHandle (+19 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "sync.Mutex"
Cohesion: 0.04
Nodes (30): lsmHook, refusedReplacements, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, checkBtrfsKeys() (+22 more)

### Community 27 - "Finder"
Cohesion: 0.17
Nodes (10): findButton, span, column(), findAll(), Finder, lowerRunes(), newFindInput(), overlay() (+2 more)

### Community 28 - "FileTarget"
Cohesion: 0.10
Nodes (14): guardPill(), newExeResolver(), Files, FileSpec, FileTarget, Policy, guardRow(), newFileTarget() (+6 more)

### Community 29 - ".step"
Cohesion: 0.15
Nodes (17): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedEntry(), placedFileStat(), placedMatches() (+9 more)

### Community 30 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.10
Nodes (15): recordTickCmd(), recordModel, listenForGuardEvents(), tickGuardStats(), listenForNetEvents(), tickNetStats(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast() (+7 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.23
Nodes (10): NewGuardUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkGuardUseCaseStartError() (+2 more)

### Community 32 - "binaryList"
Cohesion: 0.10
Nodes (13): binaryList, netGuardTab, policyEditor, policyHooks, policyTab, newBinaryList(), newPolicyTab(), verbFor() (+5 more)

### Community 33 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "fscrypt/orphans.go"
Cohesion: 0.23
Nodes (19): appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies() (+11 more)

### Community 37 - "Message"
Cohesion: 0.13
Nodes (6): checkSecret(), backend, backendOp, backendState, Message, Request

### Community 38 - "controlServer"
Cohesion: 0.08
Nodes (26): swapFile(), readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest() (+18 more)

### Community 39 - "codeedit.go"
Cohesion: 0.14
Nodes (34): class, classOf(), Edit(), leadingIndent(), moveTo(), Navigate(), New(), onlySpaces() (+26 more)

### Community 40 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (9): Hash(), IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite (+1 more)

### Community 44 - "highlight.go"
Cohesion: 0.11
Nodes (8): diffModel, editorModel, fg(), Highlighter, lexerFor(), shebangLexer(), TestTokenizeKeepsLineCount(), tokenClassOf()

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "journal"
Cohesion: 0.15
Nodes (7): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 48 - "runOfflineEdit"
Cohesion: 0.12
Nodes (17): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt() (+9 more)

### Community 49 - "fadeScroll"
Cohesion: 0.05
Nodes (28): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+20 more)

### Community 50 - "vaultfs.go"
Cohesion: 0.05
Nodes (57): entryOf(), readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), ancestorEnv() (+49 more)

### Community 51 - "cursor.go"
Cohesion: 0.09
Nodes (17): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+9 more)

### Community 52 - "proto.go"
Cohesion: 0.07
Nodes (21): Client, File, Progress, User, childEnv(), launch(), Start(), verifiedBinary() (+13 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.09
Nodes (51): systemPatterns(), systemRule(), watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp() (+43 more)

### Community 54 - "runView"
Cohesion: 0.05
Nodes (18): daemonPage, logGrid, logHook, logLine, runHost, runView, stageItem, stageState (+10 more)

### Community 55 - ".apply"
Cohesion: 0.15
Nodes (8): Apply(), Outcome, ProgressFunc, applier, ledgerOps, PlannedBinary, Store, vaultOps

### Community 57 - "mutedText"
Cohesion: 0.08
Nodes (38): filesPage, recordForm, newLogo(), newBinaryRow(), newCheck(), newCheckGroup(), newRadioGroup(), newSelect() (+30 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (55): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+47 more)

### Community 59 - "Vault"
Cohesion: 0.17
Nodes (16): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs() (+8 more)

### Community 60 - "os.FileInfo"
Cohesion: 0.24
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 61 - "Recorder"
Cohesion: 0.16
Nodes (15): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 63 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 64 - "FileEvent"
Cohesion: 0.16
Nodes (5): snapMonitor, FileEvent, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 65 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.09
Nodes (21): lineStyles, pos, viewportState, dedent(), lines(), repeatKey(), reposition(), baseStyle() (+13 more)

### Community 66 - "context.Context"
Cohesion: 0.19
Nodes (9): centered(), FitSize(), poll(), TestFitSize(), hyprClient, hyprMonitor, Rect, Size (+1 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "classifyGroup"
Cohesion: 0.20
Nodes (14): appletLine(), classifyGroup(), classifyLibs(), innermostResource(), libCovered(), newBinary(), verifyExe(), worldWritableDir() (+6 more)

### Community 69 - "liveSession"
Cohesion: 0.15
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+5 more)

### Community 70 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 71 - "batchInodeRows"
Cohesion: 0.24
Nodes (11): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+3 more)

### Community 72 - "Candidate"
Cohesion: 0.28
Nodes (13): appendSectionsAndEdit(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+5 more)

### Community 73 - "string"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "lockAndDeprovision"
Cohesion: 0.18
Nodes (8): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision(), modifiedContextWithSource()

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "shellrc.go"
Cohesion: 0.10
Nodes (53): listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), removeFacultative(), removeFacultatives() (+45 more)

### Community 78 - "update_test.go"
Cohesion: 0.05
Nodes (73): installHooks(), staged(), changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), applyUpdate(), archAssets() (+65 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.12
Nodes (8): holdHeadless(), HeadlessLine(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), GuardEvent, logBacklogDenial()

### Community 81 - "op"
Cohesion: 0.09
Nodes (25): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+17 more)

### Community 82 - "User"
Cohesion: 0.03
Nodes (83): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), maintainedApps() (+75 more)

### Community 83 - "maintainerKeeper"
Cohesion: 0.16
Nodes (12): appRunning(), libDirSig(), startMaintainers(), grantNewBinaries(), maintainedApp, maintainerKeeper, Line(), Lines() (+4 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "updatePage"
Cohesion: 0.29
Nodes (5): latestRelease, updatePage, newUpdatePage(), boolIndex(), isPreRelease()

### Community 86 - "Binary"
Cohesion: 0.20
Nodes (15): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+7 more)

### Community 87 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 88 - "dialLiveSession"
Cohesion: 0.10
Nodes (23): auditAfterEdit(), TestValidatePassword(), ValidatePassword(), dialLiveSession(), applyEdit(), editAgain(), fetchConfig(), putConfig() (+15 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "EventType"
Cohesion: 0.13
Nodes (11): BpfEvent, EventType, eventUnitTest, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), specSummary(), Cstr(), DecodeBpfEvent() (+3 more)

### Community 92 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.07
Nodes (20): foreignPinnedLink(), reasonText(), mcNameKey(), statInodeKey(), codeVouched(), isTraced(), judgePasses(), loaderEnv() (+12 more)

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
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 100 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 101 - "pinstate.go"
Cohesion: 0.29
Nodes (14): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+6 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NetGuardEvent"
Cohesion: 0.13
Nodes (12): computeGuardBinaries(), discoverInfra(), resolveGuardMode(), runTUI(), snapNetGuard, Mode, NetGuardEvent, NewNetGuardModel() (+4 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "hashFileShared"
Cohesion: 0.27
Nodes (8): HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled(), unchangedSince()

### Community 111 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 112 - "wm.go"
Cohesion: 0.29
Nodes (8): detectWM(), reachable(), Detect(), Session, hyprlandSocket(), i3RuntimeSocket(), socketOwnedBy(), Kind

### Community 113 - "recordFlow"
Cohesion: 0.15
Nodes (13): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+5 more)

### Community 115 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 116 - "Open"
Cohesion: 0.16
Nodes (20): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+12 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.04
Nodes (73): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+65 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.48
Nodes (5): copy(), is_verb(), main(), read_once(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "install/install.go"
Cohesion: 0.05
Nodes (70): releaseBuild(), Run(), installGUI(), asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureMasterKey(), fileSHA256() (+62 more)

### Community 131 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "time.Time"
Cohesion: 0.25
Nodes (5): LockedError, sweepID, sweepMemo, listEntries(), dirEntry

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "catalogRefresher"
Cohesion: 0.13
Nodes (11): drain(), newCatalogRefresher(), readEvents(), sameConfig(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), startCatalogRefresh() (+3 more)

### Community 143 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

### Community 144 - "NewDumpHook"
Cohesion: 0.19
Nodes (15): TestSetupDumpLogRefusesSymlink(), applyVerbosity(), effectiveVerbose(), preRun(), setupDumpLog(), VerboseLevel, BridgeStdLog(), NewDumpHook() (+7 more)

### Community 145 - "UI"
Cohesion: 0.10
Nodes (18): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+10 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (12): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, applySchema(), Ledger (+4 more)

### Community 147 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 148 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 149 - "changelog_test.go"
Cohesion: 0.17
Nodes (11): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+3 more)

### Community 150 - "Config"
Cohesion: 0.07
Nodes (44): newBinaryVetter(), openBinaryVetter(), makeReloadHandler(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet() (+36 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.09
Nodes (36): backingDeviceUnion(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon(), unlockPendingGroupRoots() (+28 more)

### Community 155 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 157 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+2 more)

### Community 164 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (50): TestForwardFlagSurface(), actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction (+42 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (49): fakeLink, homePage, pageDef, pageID, reel, resizeWatch, shell, sizedDialog (+41 more)

### Community 167 - "newTestTarget"
Cohesion: 0.32
Nodes (12): newTestTarget(), since(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst(), TestApplyOffStopsGuard(), TestApplyOffToWhitelistStartsGuard() (+4 more)

### Community 168 - "watchSet"
Cohesion: 0.23
Nodes (6): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet

### Community 169 - "SanitizeText"
Cohesion: 0.09
Nodes (22): multicallAdmissible(), ownerPaths(), scriptIdentity(), logRefreshChange(), TestLogRefreshChange(), formatElapsed(), orNone(), printBinary() (+14 more)

### Community 170 - "wm_test.go"
Cohesion: 0.67
Nodes (8): envOf(), serve(), shortDir(), TestDetect(), TestHyprlandFloat(), TestHyprlandFloatTimesOutWithoutWindow(), TestI3Float(), uid()

### Community 171 - "runLockdown"
Cohesion: 0.29
Nodes (8): configureDaemonLogging(), resolveConfigPath(), runLockdown(), dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "runBackend"
Cohesion: 0.36
Nodes (6): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive()

### Community 175 - "github.com/spf13/cobra.Command"
Cohesion: 0.05
Nodes (57): DrainUntilSignal(), AddHeadlessFlag(), AddServeFlags(), CheckBPFLSM(), CheckEBPF(), credentialFlagState(), ServeConfig, isLoopbackHost() (+49 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 178 - "stripTempJournal"
Cohesion: 0.32
Nodes (8): ensurePinStateFilePlaceholder(), ensurePlaceholder(), writeInPlace(), bootstrapSelfProtectPlaceholders(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace(), writeTempJournal()

### Community 179 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "KernelDev"
Cohesion: 0.09
Nodes (36): fdStat, BinaryStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), openRoot(), physicalParent() (+28 more)

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

### Community 194 - "writeWithin"
Cohesion: 0.25
Nodes (8): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "plan_test.go"
Cohesion: 0.09
Nodes (30): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), isWhiteout(), opaque(), pathByte() (+22 more)

### Community 199 - "MonitorRepository"
Cohesion: 0.25
Nodes (3): MonitorRepository, NewMonitorUseCase(), TestMonitorUseCaseStartError()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 202 - "DaemonEvent"
Cohesion: 0.06
Nodes (52): UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), confirmMasterKeyOverwrite(), drainEphemeral(), markReady(), newDaemonModel() (+44 more)

### Community 203 - "Plan"
Cohesion: 0.22
Nodes (6): dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule

### Community 204 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), noteReloadedSettings(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 205 - "RunWith"
Cohesion: 0.07
Nodes (39): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups(), restoreBackups() (+31 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - "guardModel"
Cohesion: 0.08
Nodes (24): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), appendCapped(), clipToWidth(), formatResourceBar(), formatType() (+16 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 214 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 217 - "watchPattern"
Cohesion: 0.52
Nodes (6): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern

### Community 218 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 221 - "BinaryEntry"
Cohesion: 0.03
Nodes (46): keptEntries, snapNetMonitor, BinaryEntry, NetEvent, NetEventType, NetRecord, diffBinaries(), TestDiffBinaries() (+38 more)

### Community 222 - "FileEngines"
Cohesion: 0.60
Nodes (4): FileEngines, NewFiles(), TestAddRefusesOverlap(), TestRemoveStopsGuardThenMonitor()

### Community 224 - "runWizard"
Cohesion: 0.04
Nodes (84): checkAndApply(), anyDaemonRunning(), askBunTmpdirUsers(), gatherPerUserSetup(), setupBunTmpdir(), addedResources(), applyDiffAdditions(), collectDiffAdditions() (+76 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 227 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 229 - "links_test.go"
Cohesion: 0.54
Nodes (7): newLinkView(), checkLink(), linkFixture(), TestCheckLinksAccepts(), TestCheckLinksJudgesEachTarget(), TestCheckLinksRefusesUserFiles(), TestCheckLinksSeesThePlannedTree()

### Community 234 - "Active"
Cohesion: 0.12
Nodes (15): showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestProgressReader(), cancelUI, TestEditTextMapsUICancel() (+7 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 244 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 633 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **41 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.09795863670880421 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `sync.Mutex`, `GuardInodeKey`, `time.Duration`, `controlServer`, `time.Time`, `guardUnitTest`, `DaemonEvent`, `os.File`, `hashFileShared`, `github.com/spf13/cobra.Command`, `GuardEvent`, `daemonconfig.go`, `buildOneGuard`, `fakeDaemon`, `BinaryEntry`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.019867549668874173 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `RunWith`, `Message`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09263157894736843 - nodes in this community are weakly interconnected._