# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 538 files · ~1,762,717 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7012 nodes · 31310 edges · 243 communities (203 shown, 40 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2893 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ae87dab1`
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
- go_pkg_path_filepath
- New
- protectedPage
- os.File
- plan_test.go
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- find_test.go
- guard.bpf.c
- FileTarget
- changelog.go
- BinaryStat
- golang.org/x/sys/unix.Stat_t
- Config
- TestMonitorUseCaseLifecycle
- policyEditor
- Guard
- picker
- SelectOrphans
- IntegrationSuite
- Message
- controlServer
- github.com/charmbracelet/bubbles/textarea.Model
- pinstate.go
- IntegrationSuite
- absPath
- monitorUnitTest
- BinaryEntry
- fileedit_model_test.go
- EventType
- stdio
- update_test.go
- fadeScroll
- vaultfs.go
- editView
- proto.go
- daemonconfig.go
- hashCache
- .apply
- githubRelease
- mutedText
- GuardInodeKey
- procstats.go
- audit.go
- Recorder
- inode_dev
- root.go
- installPage
- unix_connect.c
- wm.go
- IntegrationSuite
- linkView
- liveSession
- CandidateDir
- RefreshSection
- Candidate
- fcntl
- bpfstats/main.go
- SanitizeText
- codecheck/main.go
- update/update.go
- Network
- process_vm_readv.c
- loadDaemonConfig
- maintainer/apply.go
- LibraryClosure
- User
- fyne.io/fyne/v2/driver/desktop.Cursor
- journal
- runOfflineEdit
- loginServer
- GenerateConf
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- Run
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- dialLiveSession
- os.FileInfo
- Well-known bypass classes
- filevault.go
- jit_provenance.c
- inspector_probe.c
- recordops.c
- Vault
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- deploy
- binaryList
- rcUser
- recordFlow
- IntegrationSuite
- backend.go
- github.com/spf13/cobra.Command
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- ResolvePinBase
- WebSocket /ws client
- gui/session.go
- install.sh
- mc_tag.h
- Monitor
- make test-integration (rootful Docker suite)
- IntegrationSuite
- updatePage
- Open
- trace-app-libs.sh
- IntegrationSuite
- maintainerKeeper
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- debouncer
- serve.go
- NetEventType
- Prereq
- Ledger
- planUpdaters
- GlobReservations
- batchInodeRows
- sources
- fakeStore
- Finder
- IntegrationSuite
- buildOneGuard
- runNetworkMonitor
- CopyTreeWithProgress
- safeio.go
- Installation
- backend_test.go
- graphify-refresh.sh
- RawTarget
- backups.go
- newTestTarget
- NetGuardEvent
- feed.go
- shell
- layout.go
- time.Time
- sources_test.go
- WithSuspectsFrom
- Limitations
- net_tester/main.go
- IntegrationSuite
- TrustGuard
- unit.go
- tempGrantSetup
- Daemon configuration
- appTheme
- Compatibility
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- Binary
- .populateWatchInodes
- edit-protected
- headless/headless.go
- kernel-prof.sh
- ValidateTargets
- pathCache
- fileedit_symlink_test.go
- DaemonEvent
- WithBottomBar
- github.com/stretchr/testify/suite.Suite
- NewEventFanout
- Desktop GUI
- IntegrationSuite
- newBinary
- parseGuardEvent
- .eventTypeAllowed
- UI
- netinject.c
- sync.Mutex
- .awaitDaemonUpOrDead
- bpfLsmListed
- VouchedDevs
- .FilterExistingWhitelist
- TestMain
- IntegrationSuite
- downloadFile
- SecretStoreNotes
- .showSized
- GuardEvent
- supersede_probe.c
- bleedLayout
- maintainer/plan.go
- mountTable
- linkFixture
- github.com/Virgula0/app-listener/tools/codecheck
- TestCatalogSanity
- IntegrationSuite
- parsePasswd
- CheckBPFLSM
- unlockUnderGuard
- Troubleshooting
- prctl_gate
- .TestDaemon_Btrfs_CopyIoctlGate
- namedGeneralTool
- ptrace_seize.c
- .runEditorHarness
- monitor.md

## God Nodes (most connected - your core abstractions)
1. `Config` - 129 edges
2. `Guard` - 124 edges
3. `GuardInodeKey` - 103 edges
4. `shell` - 97 edges
5. `Load()` - 93 edges
6. `absPath()` - 91 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 72 edges
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

## Communities (243 total, 40 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.12
Nodes (33): atomicWriteAt(), writeAndSync(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs (+25 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (139): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload() (+131 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.13
Nodes (18): deprovisionKind, isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn() (+10 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (76): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+68 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (22): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (34): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+26 more)

### Community 11 - "New"
Cohesion: 0.10
Nodes (32): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), sectionsFromCandidates(), discordCatalogEntry(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+24 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (23): backendLink, checkGrid, configView, forwardForm, forwardView, handRadioGroup, protectedPage, recordForm (+15 more)

### Community 13 - "os.File"
Cohesion: 0.03
Nodes (85): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), setConfinedHomes() (+77 more)

### Community 14 - "plan_test.go"
Cohesion: 0.40
Nodes (11): scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+3 more)

### Community 15 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (41): recordTickCmd(), recordModel, diffModel, editorModel, formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision() (+33 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "navItem"
Cohesion: 0.04
Nodes (18): clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar, tabButton (+10 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (48): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+40 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (14): backingDeviceUnion(), fakeDaemon, Resource, BackingDevices(), GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf() (+6 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (41): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+33 more)

### Community 24 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (76): exe_refused(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+68 more)

### Community 26 - "FileTarget"
Cohesion: 0.13
Nodes (11): Target, guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec() (+3 more)

### Community 27 - "changelog.go"
Cohesion: 0.10
Nodes (44): newChangelogModel(), showChangelog(), class, classOf(), column(), dedent(), Edit(), leadingIndent() (+36 more)

### Community 28 - "BinaryStat"
Cohesion: 0.23
Nodes (11): BinaryStat, binaryStatOf(), binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches() (+3 more)

### Community 29 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.13
Nodes (21): placedWalk, treeJudge, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedEntry(), placedFileStat() (+13 more)

### Community 30 - "Config"
Cohesion: 0.04
Nodes (92): newCatalogRefresher(), makeReloadHandler(), newDaemonModel(), runTUI(), startCatalogRefresh(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced() (+84 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.15
Nodes (15): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle() (+7 more)

### Community 32 - "policyEditor"
Cohesion: 0.13
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), netGuardPill() (+4 more)

### Community 33 - "Guard"
Cohesion: 0.07
Nodes (11): exeRow, rootHandle, supersededBinary, tempMaskOp, tempRow, openRoot(), SharedPinDegraded(), TestNoteRootMountFollowsMountID() (+3 more)

### Community 34 - "picker"
Cohesion: 0.07
Nodes (32): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, lstatOf(), outsideHardlinks() (+24 more)

### Community 35 - "SelectOrphans"
Cohesion: 0.24
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (16): checkSecret(), checkTrustLine(), trustRows(), TestLiveSessionWatchAndPinger(), newRecordFlow(), backend, backendOp, backendState (+8 more)

### Community 38 - "controlServer"
Cohesion: 0.07
Nodes (29): swapFile(), readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines() (+21 more)

### Community 39 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.09
Nodes (20): lineStyles, pos, viewportState, reposition(), fg(), Highlighter, baseStyle(), gutterLabel() (+12 more)

### Community 40 - "pinstate.go"
Cohesion: 0.13
Nodes (28): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder() (+20 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.08
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (8): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 44 - "BinaryEntry"
Cohesion: 0.05
Nodes (49): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, snapEngines, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+41 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.19
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "EventType"
Cohesion: 0.07
Nodes (18): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), specFromForm() (+10 more)

### Community 48 - "update_test.go"
Cohesion: 0.08
Nodes (24): fetchReleases(), newerThanInstalled(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleaseByTag(), TestFetchReleases() (+16 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (28): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+20 more)

### Community 50 - "vaultfs.go"
Cohesion: 0.12
Nodes (23): DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead(), mustWrite() (+15 more)

### Community 51 - "editView"
Cohesion: 0.10
Nodes (8): codeEntry, editView, fileSize(), joinRel(), newCodeEntry(), newEditView(), relOf(), endedBanner()

### Community 52 - "proto.go"
Cohesion: 0.07
Nodes (22): entryOf(), Client, File, Progress, TrustItem, User, childEnv(), launch() (+14 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.11
Nodes (47): watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector() (+39 more)

### Community 54 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 55 - ".apply"
Cohesion: 0.14
Nodes (10): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier, ledgerOps, PlannedBinary (+2 more)

### Community 56 - "githubRelease"
Cohesion: 0.16
Nodes (21): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys() (+13 more)

### Community 57 - "mutedText"
Cohesion: 0.06
Nodes (46): binaryRow, filesPage, handButton, handCheck, homePage, iconTile, panel, pill (+38 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (58): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+50 more)

### Community 59 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 60 - "audit.go"
Cohesion: 0.29
Nodes (12): auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+4 more)

### Community 61 - "Recorder"
Cohesion: 0.17
Nodes (14): drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), TestMonitorTargetsNameFilesByAbsolutePath() (+6 more)

### Community 62 - "inode_dev"
Cohesion: 0.12
Nodes (23): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+15 more)

### Community 63 - "root.go"
Cohesion: 0.08
Nodes (26): installHooks(), releaseBuild(), Run(), staged(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary() (+18 more)

### Community 64 - "installPage"
Cohesion: 0.10
Nodes (14): actionGate, binSource, installAction, InstallMode, installPage, InstallRequest, Options, UninstallHooks (+6 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "linkView"
Cohesion: 0.43
Nodes (3): linkTarget(), entryKind, linkView

### Community 69 - "liveSession"
Cohesion: 0.20
Nodes (11): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents() (+3 more)

### Community 70 - "CandidateDir"
Cohesion: 0.08
Nodes (24): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+16 more)

### Community 71 - "RefreshSection"
Cohesion: 0.08
Nodes (29): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+21 more)

### Community 72 - "Candidate"
Cohesion: 0.12
Nodes (22): selectAndEditConfig(), addManualDirectories(), groupCandidates(), pickDirectories(), pickFromCandidates(), pickUsers(), TestGroupCandidates(), cleanOrphanedMetadata() (+14 more)

### Community 73 - "fcntl"
Cohesion: 0.08
Nodes (11): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), hold_forever() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "SanitizeText"
Cohesion: 0.07
Nodes (23): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), systemPatterns(), systemRule(), startTrustGuard() (+15 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "update/update.go"
Cohesion: 0.20
Nodes (21): channelOf(), deployGUI(), FetchLatest(), fetchPair(), fetchRelease(), FetchVerified(), filterChannel(), findUpdate() (+13 more)

### Community 78 - "Network"
Cohesion: 0.09
Nodes (12): NetworkGuardRepository, newExeResolver(), guardRow(), clonePolicy(), NetEngines, NetExtra, Network, NewNetwork() (+4 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), noteReloadedSettings(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 81 - "maintainer/apply.go"
Cohesion: 0.17
Nodes (11): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+3 more)

### Community 82 - "LibraryClosure"
Cohesion: 0.07
Nodes (31): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), IsScript(), ldConfInclude(), ldSoConfDirs() (+23 more)

### Community 83 - "User"
Cohesion: 0.18
Nodes (23): rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), applyRC(), bunTargets(), EnsureBunLauncherEnv(), ensureRCBlock(), HasBunLauncherEnv(), hasRCBlock() (+15 more)

### Community 84 - "fyne.io/fyne/v2/driver/desktop.Cursor"
Cohesion: 0.11
Nodes (10): buttonCursor, handCheckGroup, handSelect, handCursor(), newButtonCursor(), newCheckGroup(), newSelect(), TestButtonCursorFollowsState() (+2 more)

### Community 85 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 86 - "runOfflineEdit"
Cohesion: 0.09
Nodes (21): auditAfterEdit(), reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), TestPickEncryptedDirSingleSkipsPicker() (+13 more)

### Community 87 - "loginServer"
Cohesion: 0.15
Nodes (13): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor() (+5 more)

### Community 88 - "GenerateConf"
Cohesion: 0.07
Nodes (44): groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), LibraryBlock, AppendSectionWhitelist(), ConfSafePath(), findSectionEnd() (+36 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.05
Nodes (16): daemonPage, daemonState, logGrid, logHook, logLine, runHost, runView, journalLine() (+8 more)

### Community 92 - "Run"
Cohesion: 0.33
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.16
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.06
Nodes (28): fdStat, physicalParent(), statFD(), foreignPinnedLink(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices() (+20 more)

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

### Community 100 - "dialLiveSession"
Cohesion: 0.13
Nodes (17): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+9 more)

### Community 101 - "os.FileInfo"
Cohesion: 0.19
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "filevault.go"
Cohesion: 0.14
Nodes (26): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), isFileVaultCiphertext(), looksLikeFileVaultRecord(), newFileVaultAEAD() (+18 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 110 - "deploy"
Cohesion: 0.15
Nodes (20): applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), deploy(), downloadAndVerify(), fetchGUI() (+12 more)

### Community 111 - "binaryList"
Cohesion: 0.12
Nodes (11): binaryList, handCursors(), handCursorsOn(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog(), showForm() (+3 more)

### Community 112 - "rcUser"
Cohesion: 0.22
Nodes (18): TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), EnsureAddKeysToAgent(), EnsureSSHAgentEnv(), RemoveSSHAgentEnv(), TestEnsureSSHAgentEnvRefusesSymlinkedParent(), TestRemoveSSHAgentEnvRemovesOnlyTheCheckedFile(), rcUser() (+10 more)

### Community 113 - "recordFlow"
Cohesion: 0.17
Nodes (15): encryptBlocker(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord(), runRecordView() (+7 more)

### Community 115 - "backend.go"
Cohesion: 0.04
Nodes (111): anyDaemonRunning(), offerDaemonStart(), applyLiveRefresh(), asRoot(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled() (+103 more)

### Community 116 - "github.com/spf13/cobra.Command"
Cohesion: 0.07
Nodes (38): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+30 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "gui/session.go"
Cohesion: 0.18
Nodes (16): invokingIDs(), ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+8 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "Monitor"
Cohesion: 0.10
Nodes (5): NewMonitor(), putBtrfsLayout(), exeKey, Monitor, watchedInode

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "updatePage"
Cohesion: 0.23
Nodes (7): installedBinary, latestRelease, updatePage, readInstalled(), newUpdatePage(), boolIndex(), isPreRelease()

### Community 131 - "Open"
Cohesion: 0.10
Nodes (31): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+23 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "maintainerKeeper"
Cohesion: 0.21
Nodes (9): appRunning(), libDirSig(), grantNewBinaries(), maintainedApp, maintainerKeeper, Lines(), TestLinesKeyOnLinkAndIncludePending(), Maintainer (+1 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "debouncer"
Cohesion: 0.18
Nodes (6): drain(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), debouncer

### Community 143 - "serve.go"
Cohesion: 0.09
Nodes (31): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+23 more)

### Community 144 - "NetEventType"
Cohesion: 0.07
Nodes (22): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), BinariesSummary(), TestBinariesSummary() (+14 more)

### Community 145 - "Prereq"
Cohesion: 0.20
Nodes (6): classifySetupError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), Prereq, fakeVault

### Community 146 - "Ledger"
Cohesion: 0.09
Nodes (15): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+7 more)

### Community 147 - "planUpdaters"
Cohesion: 0.17
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 148 - "GlobReservations"
Cohesion: 0.15
Nodes (11): GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, resolveBits(), resolveChildren(), statKey() (+3 more)

### Community 149 - "batchInodeRows"
Cohesion: 0.24
Nodes (11): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+3 more)

### Community 150 - "sources"
Cohesion: 0.24
Nodes (5): TestGeneratedLibraryBlockRoundTrips(), pathByte(), pathRecord(), TestPathRecord(), sources

### Community 152 - "Finder"
Cohesion: 0.16
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.05
Nodes (62): DrainUntilSignal(), CheckEBPF(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration(), prepareDaemonStart(), reloadOnce() (+54 more)

### Community 155 - "runNetworkMonitor"
Cohesion: 0.10
Nodes (12): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString() (+4 more)

### Community 156 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "safeio.go"
Cohesion: 0.25
Nodes (13): TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), EnsureBunTmpDir(), TestEnsureBunTmpDirCreates0700(), AtomicWriteAt(), CreateExclNoFollow(), descend() (+5 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+2 more)

### Community 162 - "backups.go"
Cohesion: 0.24
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 163 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 164 - "NetGuardEvent"
Cohesion: 0.12
Nodes (15): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), snapNetGuard (+7 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (42): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+34 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (51): fakeLink, pageDef, pageID, reel, resizeWatch, shell, snapNet, sysStatus (+43 more)

### Community 167 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 168 - "time.Time"
Cohesion: 0.11
Nodes (16): LockedError, addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+8 more)

### Community 169 - "sources_test.go"
Cohesion: 0.20
Nodes (9): SystemTrusted(), systemTool(), archiveMembers(), decompressed(), newSources(), tarMembers(), TestArchiveMembersIgnoresNonArchives(), TestSourcesHoldFilesAndArchiveMembers() (+1 more)

### Community 170 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 175 - "TrustGuard"
Cohesion: 0.06
Nodes (19): lsmHook, tempGrant, trustHook, trustRows, deleteResKeys(), setSupersedeTrust(), liveTrust(), loadPinnedExeMaps() (+11 more)

### Community 176 - "unit.go"
Cohesion: 0.07
Nodes (38): TestDaemonUnitWithMetadataOutput(), dupFD(), installConfig(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit() (+30 more)

### Community 177 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 178 - "Daemon configuration"
Cohesion: 0.18
Nodes (10): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Runtime sections (+2 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

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
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "Binary"
Cohesion: 0.18
Nodes (16): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+8 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "headless/headless.go"
Cohesion: 0.29
Nodes (8): stagedBinary, stageRelease(), Check(), CheckGUI(), inspect(), buildTiny(), TestCheck(), TestCheckGUI()

### Community 199 - "ValidateTargets"
Cohesion: 0.14
Nodes (13): Section, dedupe(), existingSection(), Plan, guardedMatch(), libOfChosen(), overlaps(), refusedPath() (+5 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.07
Nodes (45): ServeConfig, UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), drainEphemeral(), runDaemon(), runDaemonUI() (+37 more)

### Community 203 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 204 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.24
Nodes (3): harnessSuite, harnessSuite, editorHarness

### Community 205 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 206 - "Desktop GUI"
Cohesion: 0.13
Nodes (15): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+7 more)

### Community 209 - "newBinary"
Cohesion: 0.29
Nodes (8): classifyLibs(), innermostResource(), libCovered(), newBinary(), worldWritableDir(), opNames(), sortedKeys(), Lib

### Community 210 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "UI"
Cohesion: 0.10
Nodes (16): confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), Active(), Confirm() (+8 more)

### Community 213 - "netinject.c"
Cohesion: 0.32
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 214 - "sync.Mutex"
Cohesion: 0.18
Nodes (7): fyneLogFilter, newActivity(), TestRingKeepsNewest(), Counts, newRing(), activity, ring

### Community 216 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 217 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 218 - ".FilterExistingWhitelist"
Cohesion: 0.09
Nodes (24): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), confSafeMatch(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob() (+16 more)

### Community 219 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 221 - "downloadFile"
Cohesion: 0.32
Nodes (5): downloadFile(), downloadReleaseFiles(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), progressReader

### Community 222 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 223 - ".showSized"
Cohesion: 0.11
Nodes (14): installUI, multiAnswer, sizedDialog, scrollable(), selectable(), fadeScrollOf(), newFadeScroll(), ask() (+6 more)

### Community 224 - "GuardEvent"
Cohesion: 0.04
Nodes (27): holdHeadless(), HeadlessLine(), forwardRequest, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+19 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 227 - "maintainer/plan.go"
Cohesion: 0.43
Nodes (6): isWhiteout(), opaque(), libPlan, op, opKind, walker

### Community 228 - "mountTable"
Cohesion: 0.19
Nodes (8): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), mountTable

### Community 229 - "linkFixture"
Cohesion: 0.43
Nodes (7): newLinkView(), checkLink(), linkFixture(), TestCheckLinksAccepts(), TestCheckLinksJudgesEachTarget(), TestCheckLinksRefusesUserFiles(), TestCheckLinksSeesThePlannedTree()

### Community 233 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 234 - "CheckBPFLSM"
Cohesion: 0.40
Nodes (5): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode()

### Community 235 - "unlockUnderGuard"
Cohesion: 0.50
Nodes (5): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), SectionScan

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 241 - "ptrace_seize.c"
Cohesion: 0.50
Nodes (6): main(), peek_buffer(), run_tracer(), run_victim(), touch(), wait_for_file()

### Community 244 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 647 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **40 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestDaemon_Btrfs_CopyIoctlGate`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `TestMain`?**
  _High betweenness centrality (0.051) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.11579651941097724 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `backups.go`, `liveSession`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.021526418786692758 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._