# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 562 files · ~1,788,138 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7303 nodes · 33378 edges · 239 communities (192 shown, 47 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2993 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `87ae00f4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- New
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- github.com/charmbracelet/bubbletea.Cmd
- User
- protectedPage
- batchInodeRows
- Network
- NetworkMonitor
- IntegrationSuite
- guard_trust.bpf.c
- fyne.io/fyne/v2/widget.Label
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- FileTarget
- guard.bpf.c
- hashCache
- install/install.go
- editView
- go_pkg_strings
- Recorder
- TestMonitorUseCaseLifecycle
- fileedit_model_test.go
- Guard
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- render.go
- pinstate_test.go
- IntegrationSuite
- absPath
- monitorUnitTest
- cursor.go
- planUpdaters
- EventType
- stdio
- update_test.go
- fadeScroll
- backend.go
- .showSized
- Client
- daemonconfig.go
- sync.Mutex
- .apply
- githubRelease
- mutedText
- GuardInodeKey
- Monitor
- IntegrationSuite
- .FilterExistingWhitelist
- inode_dev
- .filterEvent
- policyEditor
- unix_connect.c
- wm_test.go
- IntegrationSuite
- liveSession
- dialLiveSession
- Vault
- recordFlow
- journal
- fcntl
- RawTarget
- LaunchRule
- codecheck/main.go
- NetGuardEvent
- CandidateDir
- process_vm_readv.c
- os.File
- Vault
- LibraryClosure
- rcUser
- op
- bpfstats/main.go
- exe_supersede.h
- loginServer
- GenerateConf
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- Binary
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- BinaryEntry
- go_pkg_os
- Well-known bypass classes
- runLockdown
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- UnifiedDiff
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- codeedit.go
- ConfinedEntry
- installPage
- os.FileInfo
- IntegrationSuite
- Open
- releaseAssets
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- UI
- maintainerKeeper
- WebSocket /ws client
- inet_connect.c
- install.sh
- mc_tag.h
- classify
- make test-integration (rootful Docker suite)
- IntegrationSuite
- catalogRefresher
- bufio.Reader
- trace-app-libs.sh
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- vaultfs.go
- serve.go
- tempGrantSetup
- github.com/charmbracelet/bubbles/textarea.Model
- Ledger
- appTheme
- GuardEvent
- buildTrustedSet
- Daemon configuration
- fakeStore
- BinaryStat
- IntegrationSuite
- buildOneGuard
- eventUnitTest
- copyRegularAt
- IntegrationSuite
- Installation
- backend_test.go
- graphify-refresh.sh
- searchDirs
- NetEventType
- SecretStoreNotes
- .populateWatchInodes
- feed.go
- shell
- .Layout
- controlManager
- changelog_test.go
- SelectOrphans
- Limitations
- restore.go
- IntegrationSuite
- TrustGuard
- net_tester/main.go
- Plan
- deploy
- refreshAdmits
- configedit_test.go
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- RunWith
- guardLSMHooks
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- secureResources
- pathCache
- Config
- time.Time
- io.Writer
- FileEvent
- refuseOutsideHardlinks
- monitor
- IntegrationSuite
- IntegrationSuite
- lockAndDeprovision
- .eventTypeAllowed
- KernelDev
- VouchedDevs
- fileedit_symlink_test.go
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- sync.Once
- .requirePortDenied
- github.com/google/fscrypt/actions.Context
- prctl_gate
- downloadAndVerify
- .requireSecretDenial
- .runEditorHarness
- supersede_probe.c
- guardModel
- TestCleanRel
- mountTable
- github.com/Virgula0/app-listener/tools/codecheck
- Troubleshooting
- Candidate
- IntegrationSuite
- taint_victim.c
- monitor: see what touches a path
- bleedLayout
- plan_test.go
- parsePasswd

## God Nodes (most connected - your core abstractions)
1. `Config` - 146 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `shell` - 97 edges
5. `Load()` - 95 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 74 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Why not an AppArmor profile?` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (239 total, 47 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.19
Nodes (28): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+20 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (174): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity() (+166 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "New"
Cohesion: 0.08
Nodes (50): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+42 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (83): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+75 more)

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
Cohesion: 0.05
Nodes (26): addErrKind, ReplacementCheck, RegisterBusDir(), classifyAddErr(), eventMask(), addRecorder(), copySelf(), guardUnitTest (+18 more)

### Community 9 - "mkdirs"
Cohesion: 0.08
Nodes (47): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchPattern(), newWatchSet(), systemBinaryPatterns(), newFakeInotify() (+39 more)

### Community 10 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (15): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast() (+7 more)

### Community 11 - "User"
Cohesion: 0.06
Nodes (58): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), offerDaemonStart() (+50 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (20): backendLink, configView, forwardForm, forwardView, protectedPage, recordForm, recordReview, recordView (+12 more)

### Community 13 - "batchInodeRows"
Cohesion: 0.22
Nodes (11): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+3 more)

### Community 14 - "Network"
Cohesion: 0.10
Nodes (10): netGuardPill(), NetworkGuardRepository, NetworkMonitorRepository, guardRow(), clonePolicy(), NetExtra, NetPolicy, Network (+2 more)

### Community 15 - "NetworkMonitor"
Cohesion: 0.11
Nodes (9): snapNetMonitor, NetEvent, BinariesSummary(), TestBinariesSummary(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey() (+1 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.08
Nodes (9): attrOpCase, IntegrationSuite, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (61): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+53 more)

### Community 18 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (24): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+16 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (48): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+40 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (15): ownerPaths(), backingDeviceUnion(), fakeDaemon, Resource, BackingDevices(), GuardRepository, Vault, ConcurrencyLimit() (+7 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (39): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), OriginOf(), parseHash() (+31 more)

### Community 24 - "FileTarget"
Cohesion: 0.10
Nodes (12): filesPage, guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec() (+4 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (79): add_inode_to_guard(), admit_port(), btrfs_copy_gate(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial() (+71 more)

### Community 26 - "hashCache"
Cohesion: 0.21
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 27 - "install/install.go"
Cohesion: 0.05
Nodes (77): asRoot(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), fileSHA256(), installBinaryOnly() (+69 more)

### Community 28 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 29 - "go_pkg_strings"
Cohesion: 0.05
Nodes (20): parseReservedPorts(), readReservedPorts(), TestParseReservedPorts(), releaseBuild(), Run(), closeDumpFile(), Execute(), initLogger() (+12 more)

### Community 30 - "Recorder"
Cohesion: 0.13
Nodes (19): Target, newBinary(), drain(), Recorder, Target, monitorTargets(), newRecorder(), opNames() (+11 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.13
Nodes (14): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+6 more)

### Community 32 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 33 - "Guard"
Cohesion: 0.06
Nodes (16): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), exeRow, rootHandle, supersededBinary, tempMaskOp (+8 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.29
Nodes (21): busPath(), connexec(), connFD(), dial(), execInto(), failf(), get(), handshake() (+13 more)

### Community 37 - "Message"
Cohesion: 0.08
Nodes (17): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+9 more)

### Community 38 - "controlServer"
Cohesion: 0.17
Nodes (9): authSelfPeer(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), controlServer, editControlSession (+1 more)

### Community 39 - "render.go"
Cohesion: 0.12
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+6 more)

### Community 40 - "pinstate_test.go"
Cohesion: 0.29
Nodes (14): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+6 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (13): IntegrationSuite, IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase (+5 more)

### Community 44 - "cursor.go"
Cohesion: 0.08
Nodes (22): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+14 more)

### Community 45 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 46 - "EventType"
Cohesion: 0.05
Nodes (40): ensurePinStateFilePlaceholder(), ensurePlaceholder(), writeInPlace(), selfEventsFor(), bootstrapSelfProtectPlaceholders(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace() (+32 more)

### Community 48 - "update_test.go"
Cohesion: 0.10
Nodes (17): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanInstalledAcrossTagFormats() (+9 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (29): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+21 more)

### Community 51 - ".showSized"
Cohesion: 0.08
Nodes (13): binaryList, installUI, multiAnswer, resizeWatch, sizedDialog, fadeScrollOf(), newFadeScroll(), ask() (+5 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (13): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Reader (+5 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.08
Nodes (63): SecretScope, watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBlockLine(), applyConfigLine(), applyDirective() (+55 more)

### Community 54 - "sync.Mutex"
Cohesion: 0.09
Nodes (11): placedKey(), Caller, Gate, Role, Listen(), peerCred(), readFrame(), Judge (+3 more)

### Community 55 - ".apply"
Cohesion: 0.15
Nodes (8): Apply(), Outcome, ProgressFunc, applier, ledgerOps, PlannedBinary, Store, vaultOps

### Community 56 - "githubRelease"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (43): checkGrid, handButton, handCheck, homePage, latestRelease, uninstallCard, updatePage, newBinaryList() (+35 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (62): BusClient, engine, portOwner, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank() (+54 more)

### Community 59 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 60 - "IntegrationSuite"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 61 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 62 - "inode_dev"
Cohesion: 0.15
Nodes (18): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), current_is_inspector(), exe_tag_bits() (+10 more)

### Community 64 - "policyEditor"
Cohesion: 0.14
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm_test.go"
Cohesion: 0.10
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "liveSession"
Cohesion: 0.16
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+5 more)

### Community 69 - "dialLiveSession"
Cohesion: 0.07
Nodes (35): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply() (+27 more)

### Community 70 - "Vault"
Cohesion: 0.17
Nodes (14): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr(), newBoundedKeyFn(), readKey() (+6 more)

### Community 71 - "recordFlow"
Cohesion: 0.14
Nodes (14): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+6 more)

### Community 72 - "journal"
Cohesion: 0.13
Nodes (7): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 75 - "LaunchRule"
Cohesion: 0.17
Nodes (14): argvEntry, applyLaunchRule(), envName(), LaunchRule, ParseAllowEnv(), parseLaunchRule(), BuiltinLaunchRules(), canonLaunchRules() (+6 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "NetGuardEvent"
Cohesion: 0.05
Nodes (39): DrainUntilSignal(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runMonitor() (+31 more)

### Community 78 - "CandidateDir"
Cohesion: 0.06
Nodes (34): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, compileGlobNames() (+26 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.03
Nodes (81): multicallAdmissible(), scriptIdentity(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun() (+73 more)

### Community 81 - "Vault"
Cohesion: 0.20
Nodes (5): Vault, Open(), TestOpenVaultRefusesSymlink(), Unopened(), WriteFileKeepMeta()

### Community 82 - "LibraryClosure"
Cohesion: 0.10
Nodes (19): hostView(), IsScript(), LibraryClosure(), mountVouchesOwnership(), TestLibraryClosureRejectsUserWritableOriginDir(), runsNatively(), skipWithoutHostRoot(), TestHostView() (+11 more)

### Community 83 - "rcUser"
Cohesion: 0.09
Nodes (53): listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), removeFacultative(), removeFacultatives() (+45 more)

### Community 84 - "op"
Cohesion: 0.06
Nodes (44): placedWalk, treeJudge, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedEntry(), placedFileStat() (+36 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.21
Nodes (16): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+8 more)

### Community 87 - "loginServer"
Cohesion: 0.08
Nodes (21): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor() (+13 more)

### Community 88 - "GenerateConf"
Cohesion: 0.04
Nodes (83): logRefreshChange(), TestLogRefreshChange(), TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+75 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.10
Nodes (5): daemonPage, runHost, runView, captureLogs(), orDefault()

### Community 92 - "Binary"
Cohesion: 0.18
Nodes (16): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+8 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (12): foreignPinnedLink(), PumpRingbuf(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses() (+4 more)

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

### Community 100 - "BinaryEntry"
Cohesion: 0.09
Nodes (16): deferredBinary, keptEntries, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), checkBtrfsKeys(), canonicalBinaryPath(), confinedEntry() (+8 more)

### Community 101 - "go_pkg_os"
Cohesion: 0.07
Nodes (39): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+31 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "runLockdown"
Cohesion: 0.12
Nodes (18): CheckBPFLSM(), configureDaemonLogging(), confirmMasterKeyOverwrite(), resolveConfigPath(), runBPFCheck(), runGenKey(), runLockdown(), runOneShotMode() (+10 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "UnifiedDiff"
Cohesion: 0.09
Nodes (22): applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), refreshCatalogConfig(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr(), confirmUnit(), TestUnifiedDiff() (+14 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "codeedit.go"
Cohesion: 0.09
Nodes (43): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), moveTo(), Navigate() (+35 more)

### Community 111 - "ConfinedEntry"
Cohesion: 0.12
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 112 - "installPage"
Cohesion: 0.07
Nodes (28): anyDaemonRunning(), actionGate, binSource, installAction, installedBinary, InstallMode, installPage, InstallRequest (+20 more)

### Community 113 - "os.FileInfo"
Cohesion: 0.16
Nodes (7): entryOf(), migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), FileModeToUnixMode()

### Community 115 - "Open"
Cohesion: 0.14
Nodes (22): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+14 more)

### Community 116 - "releaseAssets"
Cohesion: 0.21
Nodes (13): archAssets(), assetsFor(), checkUpdatePreconditions(), guiAssetsFor(), resolveAssets(), sanityCheckBinary(), TestAssetsFor(), TestGUIAssetsFor() (+5 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "UI"
Cohesion: 0.08
Nodes (22): nextConfig(), showChangelog(), channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI() (+14 more)

### Community 121 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "classify"
Cohesion: 0.15
Nodes (15): call(), TestClassify(), TestClassifyOpenCallsExact(), TestReadFrame(), arg(), attributesOf(), classify(), classifyCollection() (+7 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "catalogRefresher"
Cohesion: 0.14
Nodes (9): drain(), readEvents(), sameConfig(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), catalogRefresher (+1 more)

### Community 131 - "bufio.Reader"
Cohesion: 0.12
Nodes (17): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+9 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "vaultfs.go"
Cohesion: 0.06
Nodes (47): TestDaemonUnitWithMetadataOutput(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt(), writeAtNoFollow(), isInstallerSSHAgentUnit() (+39 more)

### Community 143 - "serve.go"
Cohesion: 0.09
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 144 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 145 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.08
Nodes (34): findButton, pos, span, column(), reposition(), findAll(), Finder, lowerRunes() (+26 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (26): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+18 more)

### Community 148 - "GuardEvent"
Cohesion: 0.05
Nodes (19): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+11 more)

### Community 149 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 150 - "Daemon configuration"
Cohesion: 0.17
Nodes (12): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Port sections, Reloading (+4 more)

### Community 152 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.07
Nodes (38): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), displayEntries() (+30 more)

### Community 155 - "eventUnitTest"
Cohesion: 0.18
Nodes (7): BpfEvent, eventUnitTest, NetRecord, Cstr(), DecodeBpfEvent(), FormatAddr(), Ntohs()

### Community 156 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 162 - "NetEventType"
Cohesion: 0.14
Nodes (11): NetEventType, eventsOrAll(), NetEventTypes(), ParseNetEventType(), WithSuspectsFrom(), eventsetSummary(), modeLabel(), NewNetGuard() (+3 more)

### Community 163 - "SecretStoreNotes"
Cohesion: 0.25
Nodes (10): isInsidePath(), PathCovered(), TestPathCovered(), dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner() (+2 more)

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (50): TestForwardFlagSurface(), actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction (+42 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (59): fakeLink, pageDef, pageID, reel, shell, snapEngines, snapNet, newLogo() (+51 more)

### Community 167 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 168 - "controlManager"
Cohesion: 0.24
Nodes (4): newControlManager(), startControlManager(), configEditor, controlManager

### Community 169 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 170 - "SelectOrphans"
Cohesion: 0.30
Nodes (13): appProtectorSet(), isAppProtector(), orphanedPolicies(), orphanedProtectors(), SelectOrphans(), appProtector(), foreignProtector(), policy() (+5 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "restore.go"
Cohesion: 0.21
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 175 - "TrustGuard"
Cohesion: 0.05
Nodes (26): statFunc, tempGrant, trustRows, ensureBtrfsLayout(), setSupersedeTrust(), liveTrust(), cStr(), TestPlanTrustedDirs_DuplicateWriterAddsNoEmptyLoader() (+18 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "Plan"
Cohesion: 0.22
Nodes (6): dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule

### Community 178 - "deploy"
Cohesion: 0.20
Nodes (14): deploy(), deployGUI(), fetchGUI(), hasAsset(), InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run() (+6 more)

### Community 179 - "refreshAdmits"
Cohesion: 0.25
Nodes (8): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), systemPatterns(), systemRule(), BinaryRule, materialized(), Extension

### Community 180 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

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
Cohesion: 0.05
Nodes (43): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+35 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "RunWith"
Cohesion: 0.29
Nodes (8): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), Mode, Request

### Community 195 - "guardLSMHooks"
Cohesion: 0.25
Nodes (5): lsmHook, trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), TestSocketHooksAreRequired()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "secureResources"
Cohesion: 0.19
Nodes (10): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse(), harnessSuite (+2 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "Config"
Cohesion: 0.04
Nodes (86): CheckEBPF(), newBinaryVetter(), openBinaryVetter(), newCatalogRefresher(), buildGuards(), loadDaemonConfig(), lockOneRoot(), lockRootRecovering() (+78 more)

### Community 202 - "time.Time"
Cohesion: 0.08
Nodes (40): LockedError, ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI() (+32 more)

### Community 203 - "io.Writer"
Cohesion: 0.12
Nodes (11): formatElapsed(), orNone(), printBinary(), printRecording(), TestProgressReader(), fyneLogFilter, bottomBar, newBottomBar() (+3 more)

### Community 204 - "FileEvent"
Cohesion: 0.12
Nodes (6): snapMonitor, FileEvent, MonitorRepository, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 205 - "refuseOutsideHardlinks"
Cohesion: 0.28
Nodes (8): lstatOf(), outsideHardlinks(), refuseOutsideHardlinks(), TestRefuseOutsideHardlinks_FileRoot(), TestRefuseOutsideHardlinks_InTreeLinksAccepted(), TestRefuseOutsideHardlinks_MissingRootFails(), TestRefuseOutsideHardlinks_OutsideLinkRefused(), writeFile()

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (30): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+22 more)

### Community 210 - "lockAndDeprovision"
Cohesion: 0.33
Nodes (6): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "KernelDev"
Cohesion: 0.12
Nodes (19): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), reasonText() (+11 more)

### Community 213 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 214 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 219 - "github.com/google/fscrypt/actions.Context"
Cohesion: 0.40
Nodes (4): modifiedContextWithSource(), CleanOrphans(), listPolicies(), listProtectors()

### Community 220 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 221 - "downloadAndVerify"
Cohesion: 0.15
Nodes (18): installHooks(), staged(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), FetchLatest(), fetchPair(), fetchRelease() (+10 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "guardModel"
Cohesion: 0.05
Nodes (39): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+31 more)

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+2 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 238 - "Candidate"
Cohesion: 0.13
Nodes (20): appendSectionsAndEdit(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickFromCandidates(), runConfigEditor(), sectionsFromCandidates() (+12 more)

### Community 249 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 256 - "plan_test.go"
Cohesion: 0.07
Nodes (37): rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), systemTool(), isWhiteout(), opaque(), pathByte(), pathRecord() (+29 more)

### Community 264 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 665 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **47 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `.requirePortDenied`, `IntegrationSuite`, `go_pkg_strings`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.016933885786344804 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `time.Duration`, `GuardInodeKey`, `BinaryEntry`, `guardUnitTest`, `Config`, `time.Time`, `batchInodeRows`, `TrustGuard`, `os.File`, `ConfinedEntry`, `GuardEvent`, `daemonconfig.go`, `Resource`, `loginServer`, `sync.Mutex`, `buildOneGuard`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `How it works`, `Installation`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `New` be split into smaller, more focused modules?**
  _Cohesion score 0.07547169811320754 - nodes in this community are weakly interconnected._