# Graph Report - app-listener  (2026-10-27)

## Corpus Check
- 525 files · ~1,741,490 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6810 nodes · 30146 edges · 244 communities (201 shown, 43 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2772 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `243546af`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- TempBinary
- github.com/spf13/cobra.Command
- protectedPage
- find_test.go
- github.com/charmbracelet/bubbletea.Model
- installPage
- IntegrationSuite
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- Guard
- guard.bpf.c
- TrustGuard
- Finder
- FileTarget
- BinaryStat
- filevault_test.go
- NetGuardEvent
- policyEditor
- changelog.go
- picker
- CandidateDir
- IntegrationSuite
- Message
- controlServer
- runDaemon
- NetworkMonitor
- IntegrationSuite
- absPath
- monitorUnitTest
- runNetworkGuard
- fileedit_model_test.go
- User
- os.File
- fadeScroll
- vaultfs.go
- cursor.go
- Client
- watchGroup
- tabs.go
- .apply
- KernelDev
- mutedText
- GuardInodeKey
- runWizard
- expandPlaceholders
- Recorder
- exe_supersede.h
- NetEventType
- plan_test.go
- removeFacultatives
- wm.go
- IntegrationSuite
- .FilterExistingWhitelist
- liveSession
- editView
- journal
- SuperblockDev
- fcntl
- bpfstats/main.go
- Network
- codecheck/main.go
- rcUser
- deploy
- process_vm_readv.c
- GuardEvent
- run.go
- SelectOrphans
- maintainerKeeper
- connFD
- github.com/charmbracelet/bubbles/textarea.Model
- Binary
- planUpdaters
- runOfflineEdit
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
- sync.Mutex
- jit_provenance.c
- inspector_probe.c
- recordops.c
- NewEventFanout
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- batchInodeRows
- mountTable
- ListUsers
- recordFlow
- IntegrationSuite
- NewDumpHook
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
- UnifiedDiff
- LibraryClosure
- trace-app-libs.sh
- IntegrationSuite
- RunWith
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- Config
- serveWebSocket
- debouncer
- UI
- Ledger
- Serve
- validateConfigText
- readGrantRequest
- buildTrustedSet
- fakeStore
- inodeChain
- IntegrationSuite
- buildOneGuard
- safeio.go
- auditGroupCoverage
- runLockdown
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- IntegrationSuite
- .populateWatchInodes
- feed.go
- shell
- newTestTarget
- watchSet
- io.Writer
- options.go
- .showSized
- Limitations
- Candidate
- IntegrationSuite
- uninstallAll
- netinject.c
- copyRegularAt
- launchKey
- TestMain
- Compatibility
- How it works
- listFacultatives
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- githubRelease
- RawTarget
- edit-protected
- AddServeFlags
- kernel-prof.sh
- update_test.go
- pathCache
- hashFileShared
- time.Time
- Plan
- editControlSession
- find
- Desktop GUI
- IntegrationSuite
- netGuardModel
- downloadFile
- .eventTypeAllowed
- bleedLayout
- WithBottomBar
- refuseOutsideHardlinks
- .awaitDaemonUpOrDead
- TestCleanRel
- ReadInstallCert
- tempGrantSetup
- ParseServeFlags
- shQuote
- ptrace_race.c
- prctl_gate
- compileGlobNames
- RunCmd
- supersede_probe.c
- WithSuspectsFrom
- ParseEventsFlag
- runMonitor
- formatEntry
- github.com/Virgula0/app-listener/tools/codecheck
- Vault
- parseKernelRelease
- taint_victim.c
- copyXattrs
- Execute
- Troubleshooting
- bpfLsmListed
- golang.org/x/sys/unix.Timespec
- fileID
- MulticallError
- monitor.md

## God Nodes (most connected - your core abstractions)
1. `Guard` - 122 edges
2. `Config` - 118 edges
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

## Communities (244 total, 43 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.08
Nodes (21): tempJournalRow, Progress, User, invalidPassword, recordTick, linkedFile, launchRule, trustEvent (+13 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (164): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+156 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): Hash(), configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.09
Nodes (29): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), deprovisionKind, harnessSuite (+21 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (74): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+66 more)

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
Nodes (24): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+16 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (37): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+29 more)

### Community 10 - "TempBinary"
Cohesion: 0.08
Nodes (23): forwardRequest, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), planTempBlock(), ResolveTempBinary() (+15 more)

### Community 11 - "github.com/spf13/cobra.Command"
Cohesion: 0.22
Nodes (12): TestSetupDumpLogRefusesSymlink(), asRoot(), runInstall(), TestMovedFlagError(), closeDumpFile(), effectiveVerbose(), headlessRequested(), movedFlagError() (+4 more)

### Community 12 - "protectedPage"
Cohesion: 0.06
Nodes (19): backendLink, configView, forwardForm, forwardView, protectedPage, recordForm, recordReview, recordView (+11 more)

### Community 13 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 14 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.06
Nodes (15): recordTickCmd(), recordModel, diffModel, editorModel, Highlighter, renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented() (+7 more)

### Community 15 - "installPage"
Cohesion: 0.07
Nodes (18): anyDaemonRunning(), actionGate, binSource, handCheck, installAction, installedBinary, installPage, uninstallCard (+10 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "runView"
Cohesion: 0.05
Nodes (15): daemonPage, logGrid, logHook, logLine, runHost, runView, stageItem, stageState (+7 more)

### Community 19 - "fileEditModel"
Cohesion: 0.07
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (11): ownerPaths(), fakeDaemon, Resource, GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots() (+3 more)

### Community 23 - "auth_test.go"
Cohesion: 0.10
Nodes (46): CheckPassword(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+38 more)

### Community 24 - "Guard"
Cohesion: 0.05
Nodes (29): binaryVerifyState, deferredBinary, keptEntries, rootHandle, BinaryEntry, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+21 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "TrustGuard"
Cohesion: 0.05
Nodes (26): exeRow, lsmHook, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, ensureBtrfsLayout() (+18 more)

### Community 27 - "Finder"
Cohesion: 0.15
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 28 - "FileTarget"
Cohesion: 0.11
Nodes (13): Target, FileEngines, Files, FileSpec, FileTarget, guardRow(), newFileTarget(), NormalizeSpec() (+5 more)

### Community 29 - "BinaryStat"
Cohesion: 0.11
Nodes (23): BinaryStat, placedWalk, binaryStatOf(), StatBinary(), StatBinaryFile(), confineBelow(), inUserHome(), durableFS() (+15 more)

### Community 30 - "filevault_test.go"
Cohesion: 0.09
Nodes (49): decryptDirectories(), TestDecryptDirectoriesEmpty(), classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord() (+41 more)

### Community 31 - "NetGuardEvent"
Cohesion: 0.11
Nodes (15): snapNetGuard, NetGuardEvent, NetworkGuardRepository, NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle() (+7 more)

### Community 32 - "policyEditor"
Cohesion: 0.13
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+4 more)

### Community 33 - "changelog.go"
Cohesion: 0.08
Nodes (48): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), class, classOf(), column() (+40 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (31): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+23 more)

### Community 35 - "CandidateDir"
Cohesion: 0.16
Nodes (18): bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), bunLaunchersOffered(), Discover(), DiscoverForUsers(), DiscoverSystem(), CandidateDir (+10 more)

### Community 37 - "Message"
Cohesion: 0.11
Nodes (8): checkSecret(), File, backend, backendOp, backendState, Message, Request, RequireDaemonStopped()

### Community 39 - "runDaemon"
Cohesion: 0.07
Nodes (37): CheckEBPF(), awaitStartupOrSignal(), backingDeviceUnion(), buildGuards(), catchLifecycleSignals(), lockOneRoot(), lockRootRecovering(), newPinGeneration() (+29 more)

### Community 40 - "NetworkMonitor"
Cohesion: 0.14
Nodes (7): BinariesSummary(), TestBinariesSummary(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 41 - "IntegrationSuite"
Cohesion: 0.08
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "runNetworkGuard"
Cohesion: 0.13
Nodes (20): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), runHeadless() (+12 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "User"
Cohesion: 0.18
Nodes (13): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), pickUsers(), removeBunLauncherEnv(), revertBunLaunchers(), globBuilder (+5 more)

### Community 48 - "os.File"
Cohesion: 0.05
Nodes (70): multicallAdmissible(), scriptIdentity(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun() (+62 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+19 more)

### Community 50 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 51 - "cursor.go"
Cohesion: 0.07
Nodes (24): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+16 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (12): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Reader (+4 more)

### Community 53 - "watchGroup"
Cohesion: 0.10
Nodes (40): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+32 more)

### Community 54 - "tabs.go"
Cohesion: 0.03
Nodes (26): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel (+18 more)

### Community 55 - ".apply"
Cohesion: 0.14
Nodes (10): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier, ledgerOps, PlannedBinary (+2 more)

### Community 56 - "KernelDev"
Cohesion: 0.11
Nodes (23): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+15 more)

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (41): checkGrid, filesPage, homePage, latestRelease, updatePage, newBinaryList(), newBinaryRow(), selectable() (+33 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (51): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+43 more)

### Community 59 - "runWizard"
Cohesion: 0.11
Nodes (26): askBunTmpdirUsers(), gatherPerUserSetup(), addedResources(), applyDiffAdditions(), TestAddedResourcesSkipsAlreadyGuardedSSH(), TestUncoveredCandidates(), uncoveredCandidates(), checkInstalledMatchesSource() (+18 more)

### Community 60 - "expandPlaceholders"
Cohesion: 0.13
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 61 - "Recorder"
Cohesion: 0.13
Nodes (17): drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), record() (+9 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.15
Nodes (19): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+11 more)

### Community 63 - "NetEventType"
Cohesion: 0.08
Nodes (15): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), FormatAddr(), NetEventTypes(), Ntohs() (+7 more)

### Community 64 - "plan_test.go"
Cohesion: 0.08
Nodes (35): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), checkLink(), linkFixture(), TestCheckLinksAccepts() (+27 more)

### Community 65 - "removeFacultatives"
Cohesion: 0.10
Nodes (30): TestDaemonUnitWithMetadataOutput(), dupFD(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit(), TestInstallSSHAgentRefusesSymlinkedUnitDir() (+22 more)

### Community 66 - "wm.go"
Cohesion: 0.07
Nodes (40): invokingIDs(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 69 - "liveSession"
Cohesion: 0.18
Nodes (12): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+4 more)

### Community 70 - "editView"
Cohesion: 0.12
Nodes (7): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 71 - "journal"
Cohesion: 0.13
Nodes (7): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 72 - "SuperblockDev"
Cohesion: 0.33
Nodes (6): FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 73 - "fcntl"
Cohesion: 0.09
Nodes (11): main(), report(), denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.09
Nodes (31): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), ensureMasterKey(), checkKeyLen(), GenerateMasterKey() (+23 more)

### Community 75 - "Network"
Cohesion: 0.16
Nodes (6): netGuardPill(), clonePolicy(), NetExtra, NetPolicy, Network, verifyHashes()

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "rcUser"
Cohesion: 0.12
Nodes (41): setupBunTmpdir(), rcTarget, BunTmpDir(), TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent() (+33 more)

### Community 78 - "deploy"
Cohesion: 0.11
Nodes (33): installHooks(), staged(), applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), deployGUI() (+25 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.19
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.04
Nodes (24): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials() (+16 more)

### Community 81 - "run.go"
Cohesion: 0.08
Nodes (26): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+18 more)

### Community 82 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 83 - "maintainerKeeper"
Cohesion: 0.20
Nodes (10): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp, maintainerKeeper (+2 more)

### Community 84 - "connFD"
Cohesion: 0.29
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (18): lineStyles, pos, viewportState, reposition(), baseStyle(), gutterLabel(), lineLen(), markStyle() (+10 more)

### Community 86 - "Binary"
Cohesion: 0.14
Nodes (20): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+12 more)

### Community 87 - "planUpdaters"
Cohesion: 0.11
Nodes (22): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+14 more)

### Community 88 - "runOfflineEdit"
Cohesion: 0.10
Nodes (22): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive() (+14 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "EventType"
Cohesion: 0.07
Nodes (18): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), specFromForm() (+10 more)

### Community 92 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (14): foreignPinnedLink(), PumpRingbuf(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs() (+6 more)

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
Nodes (14): procStats, statsSample, formatAgo(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "ConfinedEntry"
Cohesion: 0.11
Nodes (27): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache, ConfinedEntry() (+19 more)

### Community 101 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "sync.Mutex"
Cohesion: 0.11
Nodes (11): newControlManager(), startControlManager(), startControlServer(), configEditor, controlManager, newActivity(), TestRingKeepsNewest(), Counts (+3 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (14): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+6 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 110 - "batchInodeRows"
Cohesion: 0.18
Nodes (13): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+5 more)

### Community 111 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 112 - "ListUsers"
Cohesion: 0.12
Nodes (19): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), readPutContent() (+11 more)

### Community 113 - "recordFlow"
Cohesion: 0.16
Nodes (15): encryptBlocker(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord(), runRecordView() (+7 more)

### Community 115 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

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
Nodes (77): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+69 more)

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
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "UnifiedDiff"
Cohesion: 0.12
Nodes (20): applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), refreshCatalogConfig(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr(), confirmUnit(), TestUnifiedDiff() (+12 more)

### Community 131 - "LibraryClosure"
Cohesion: 0.07
Nodes (28): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+20 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "RunWith"
Cohesion: 0.14
Nodes (17): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), stagedBinary, Mode, Request (+9 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "Config"
Cohesion: 0.07
Nodes (39): ServeConfig, newBinaryVetter(), openBinaryVetter(), drain(), newCatalogRefresher(), readEvents(), swapFile(), makeReloadHandler() (+31 more)

### Community 143 - "serveWebSocket"
Cohesion: 0.15
Nodes (12): configureWebSocket(), decodeEndsAtEOF(), parseResize(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestParseResize(), TestSameOrigin() (+4 more)

### Community 144 - "debouncer"
Cohesion: 0.33
Nodes (3): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), debouncer

### Community 145 - "UI"
Cohesion: 0.08
Nodes (22): nextConfig(), showChangelog(), channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI() (+14 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "Serve"
Cohesion: 0.17
Nodes (17): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+9 more)

### Community 148 - "validateConfigText"
Cohesion: 0.11
Nodes (15): collectFilesystemPrereqs(), discordCatalogEntry(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips(), TestSectionsFromCandidatesSkipsMissingWatchSubPaths() (+7 more)

### Community 149 - "readGrantRequest"
Cohesion: 0.17
Nodes (13): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+5 more)

### Community 150 - "buildTrustedSet"
Cohesion: 0.19
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 152 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.07
Nodes (44): DrainUntilSignal(), buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions(), selfProtectSpecs(), startSelfGuard() (+36 more)

### Community 155 - "safeio.go"
Cohesion: 0.21
Nodes (15): readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), removeCheckedFile(), AtomicWriteAt() (+7 more)

### Community 156 - "auditGroupCoverage"
Cohesion: 0.11
Nodes (16): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+8 more)

### Community 157 - "runLockdown"
Cohesion: 0.15
Nodes (14): configureDaemonLogging(), loadDaemonConfig(), resolveConfigPath(), runLockdown(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), dirIsRootOwnedSafe() (+6 more)

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
Cohesion: 0.06
Nodes (38): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+30 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (47): fakeLink, pageDef, pageID, reel, shell, snapEngines, snapNet, newLogo() (+39 more)

### Community 167 - "newTestTarget"
Cohesion: 0.22
Nodes (17): newExeResolver(), NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 168 - "watchSet"
Cohesion: 0.10
Nodes (20): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), fakeInotify, watchDir (+12 more)

### Community 169 - "io.Writer"
Cohesion: 0.15
Nodes (14): formatElapsed(), orNone(), printBinary(), printRecording(), candidates(), confirmable(), confirmAll(), describe() (+6 more)

### Community 170 - "options.go"
Cohesion: 0.36
Nodes (9): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+1 more)

### Community 171 - ".showSized"
Cohesion: 0.08
Nodes (15): binaryList, installUI, multiAnswer, resizeWatch, sizedDialog, BinaryCheck, scrollable(), fadeScrollOf() (+7 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Candidate"
Cohesion: 0.26
Nodes (14): appendSectionsAndEdit(), collectDiffAdditions(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories() (+6 more)

### Community 175 - "uninstallAll"
Cohesion: 0.20
Nodes (12): releaseBuild(), Run(), confirmKeyDeletion(), run(), runUninstall(), RunWith(), removeKeyAndEmptyDir(), removeMasterKey() (+4 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "copyRegularAt"
Cohesion: 0.36
Nodes (9): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), fchownFromInfo(), futimesFromInfo(), readDirNames(), readlinkFD() (+1 more)

### Community 179 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "listFacultatives"
Cohesion: 0.24
Nodes (11): installFacultatives(), installSSHAgentFacultative(), listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH() (+3 more)

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

### Community 194 - "githubRelease"
Cohesion: 0.18
Nodes (22): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), rel(), tags(), TestChangelogText(), TestCumulativeChangelog() (+14 more)

### Community 195 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "AddServeFlags"
Cohesion: 0.43
Nodes (7): AddHeadlessFlag(), AddServeFlags(), init(), init(), init(), init(), init()

### Community 199 - "update_test.go"
Cohesion: 0.09
Nodes (22): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+14 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "hashFileShared"
Cohesion: 0.29
Nodes (8): HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled(), unchangedSince()

### Community 202 - "time.Time"
Cohesion: 0.07
Nodes (37): LockedError, UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable() (+29 more)

### Community 203 - "Plan"
Cohesion: 0.20
Nodes (6): dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule

### Community 204 - "editControlSession"
Cohesion: 0.36
Nodes (4): readClientLines(), streamEvents(), TestEditControlSessionEndIdempotent(), editControlSession

### Community 205 - "find"
Cohesion: 0.26
Nodes (14): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), Delete(), Find(), find(), Backup (+6 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - "netGuardModel"
Cohesion: 0.04
Nodes (46): runNetworkMonitor(), runTUI(), Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace() (+38 more)

### Community 210 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 214 - "refuseOutsideHardlinks"
Cohesion: 0.28
Nodes (8): lstatOf(), outsideHardlinks(), refuseOutsideHardlinks(), TestRefuseOutsideHardlinks_FileRoot(), TestRefuseOutsideHardlinks_InTreeLinksAccepted(), TestRefuseOutsideHardlinks_MissingRootFails(), TestRefuseOutsideHardlinks_OutsideLinkRefused(), writeFile()

### Community 217 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 218 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 219 - "ParseServeFlags"
Cohesion: 0.33
Nodes (7): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials()

### Community 220 - "shQuote"
Cohesion: 0.19
Nodes (4): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite

### Community 221 - "ptrace_race.c"
Cohesion: 0.52
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 223 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 224 - "RunCmd"
Cohesion: 0.05
Nodes (55): offerDaemonStart(), runDiffCatalog(), installGUI(), buildBinaryIfNeeded(), certifyInstall(), ensureInstalledBinary(), installBinaryOnly(), markInstallComplete() (+47 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 226 - "WithSuspectsFrom"
Cohesion: 0.43
Nodes (3): WithSuspectsFrom(), netguardHarness, Option

### Community 227 - "ParseEventsFlag"
Cohesion: 0.33
Nodes (4): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents()

### Community 228 - "runMonitor"
Cohesion: 0.47
Nodes (5): MakeDisplayPaths(), runHeadless(), runMonitor(), runTUI(), NewModel()

### Community 229 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 232 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 234 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 235 - "Execute"
Cohesion: 0.67
Nodes (3): Execute(), initLogger(), main()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "bpfLsmListed"
Cohesion: 0.40
Nodes (6): bpfLsmListed(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 244 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 632 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **43 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.07973342578681768 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `find`, `Message`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.018782249742002065 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._