# Graph Report - app-listener  (2026-11-06)

## Corpus Check
- 567 files · ~1,806,944 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7450 nodes · 34072 edges · 260 communities (211 shown, 49 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3042 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `229d6fe3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- op
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- .dispatch
- runWizard
- FileTarget
- ConfinedEntry
- buildOneGuard
- DaemonEvent
- IntegrationSuite
- guard_trust.bpf.c
- tabs.go
- fileEditModel
- networkmonitor.bpf.c
- newFakeVault
- Resource
- auth.go
- plan_test.go
- guard.bpf.c
- highlight_test.go
- BinaryEntry
- gui/session.go
- fileedit_model_test.go
- KernelDev
- DeployInstalledBinary
- GuardEvent
- codeedit.go
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- editView
- runMonitor
- filevault.go
- daemon.conf reference
- FindDaemonProcesses
- monitorUnitTest
- pinstate.go
- fakeGuardRepo
- updatePage
- stdio
- os.File
- fadeScroll
- LibraryClosure
- protectedPage
- Client
- Finder
- session
- .apply
- backend_test.go
- UI
- engine
- NetEventType
- TrustGuard
- serve.go
- newTestTarget
- golang.org/x/sys/unix.Stat_t
- compileGlobNames
- unix_connect.c
- main_test.go
- IntegrationSuite
- Open
- parseDaemonEvents
- Vault
- linkView
- planUpdaters
- fcntl
- binaryList
- Guard
- codecheck/main.go
- vaultfs.go
- liveSession
- process_vm_readv.c
- GuardInodeKey
- controlServer
- runGenKey
- rcUser
- inodeChain
- bpfstats/main.go
- inode_dev
- loginServer
- installPage
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- buildTrustedSet
- batchInodeRows
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- run.go
- sources
- Well-known bypass classes
- SanitizeText
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- ListUsers
- Build & sign release assets (reusable workflow)
- Candidate
- github.com/charmbracelet/bubbles/textarea.Model
- serverAuth
- netGuardModel
- recordFlow
- IntegrationSuite
- classify
- Recorder
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- time.Time
- IntegrationSuite
- WebSocket /ws client
- User
- inet_connect.c
- install.sh
- procstats.go
- EventType
- make test-integration (rootful Docker suite)
- SelectOrphans
- fakeDaemon
- github.com/spf13/cobra.Command
- trace-app-libs.sh
- IntegrationSuite
- .populateWatchInodes
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- validateConfigText
- dialLiveSession
- Ledger
- MonitorRepository
- maintainerKeeper
- btrfsLayoutFrom
- runNetworkMonitor
- RefuseIdentityLines
- fakeStore
- watchSet
- IntegrationSuite
- Monitor
- debouncer
- io/fs.FileInfo
- .admitGated
- Installation
- wire.go
- graphify-refresh.sh
- layout.go
- update_test.go
- reloadOnce
- absPath
- feed.go
- shell
- safeio.go
- github.com/charmbracelet/bubbletea.Cmd
- changelog.go
- vaultfs_test.go
- configedit_test.go
- Limitations
- cursor.go
- PlanDaemonUnitUpdate
- net_tester/main.go
- daemon: permanent protection with encryption at rest
- Daemon configuration
- fileedit_symlink_test.go
- IntegrationSuite
- NetGuardEvent
- sync.Mutex
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- catalogRefresher
- Self-updating apps and the catalog
- .netIntegrityContainer
- WithBottomBar
- decryptStep
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- RunWith
- pathCache
- os.FileInfo
- Gate
- mc_tag.h
- IntegrationSuite
- runView
- monitor
- IntegrationSuite
- IntegrationSuite
- RefreshSection
- maintainer/plan.go
- UnifiedDiff
- watchPattern
- sources_test.go
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- ReadInstallCert
- logGrid
- find_test.go
- launchKey
- uninstallAll
- ValidateTargets
- bpfLsmListed
- runNetworkGuard
- supersede_probe.c
- app-listener documentation
- prctl_gate
- loadDaemonConfig
- VouchedDevs
- github.com/Virgula0/app-listener/tools/codecheck
- SecretMatch
- taint_victim.c
- network monitor: see what a program talks to
- mutedText
- filevault_test.go
- Troubleshooting
- mountTable
- GenerateConf
- Vault
- IntegrationSuite
- parseKernelRelease
- connectionGone
- golang.org/x/sys/unix.Timespec
- fileID
- monitor: see what touches a path
- MulticallError
- netModel
- SecretStoreNotes
- expandPlaceholders
- .eventTypeAllowed
- bleedLayout
- Config
- TestCleanRel
- .keyringObject
- hashFileShared
- .runEditorHarness

## God Nodes (most connected - your core abstractions)
1. `Config` - 147 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 113 edges
4. `Load()` - 97 edges
5. `shell` - 97 edges
6. `absPath()` - 94 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 75 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Modes` --references--> `monitor()`  [INFERRED]
  docs/README.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (260 total, 49 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.09
Nodes (21): tempJournalRow, Progress, User, invalidPassword, recordTick, linkedFile, argvMatch, launchRule (+13 more)

### Community 1 - "testing.T"
Cohesion: 0.01
Nodes (192): plantBun(), TestOpenBunTmp_KeepsVettedDir(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity() (+184 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote(), watchPathsInConfig()

### Community 3 - "op"
Cohesion: 0.17
Nodes (12): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+4 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (84): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+76 more)

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
Nodes (24): addErrKind, ReplacementCheck, classifyAddErr(), addRecorder(), copySelf(), guardUnitTest, runTool(), walkEntries() (+16 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (38): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetHomeDepartureTriggersRefresh(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+30 more)

### Community 10 - ".dispatch"
Cohesion: 0.08
Nodes (19): bpfGuardEvent, BusClient, AttestBusClient(), busClientKey(), execStamp(), imageOf(), pidfdPid(), procIdentity() (+11 more)

### Community 11 - "runWizard"
Cohesion: 0.07
Nodes (45): askBunTmpdirUsers(), gatherPerUserSetup(), applyDiffAdditions(), offerDaemonStart(), runDiffCatalog(), installFacultatives(), installSSHAgentFacultative(), applyLiveRefresh() (+37 more)

### Community 12 - "FileTarget"
Cohesion: 0.12
Nodes (10): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+2 more)

### Community 13 - "ConfinedEntry"
Cohesion: 0.23
Nodes (22): ConfinedEntry(), HashStats(), newVerifyState(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash() (+14 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.06
Nodes (43): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), displayEntries() (+35 more)

### Community 15 - "DaemonEvent"
Cohesion: 0.09
Nodes (37): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+29 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (64): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+56 more)

### Community 18 - "tabs.go"
Cohesion: 0.03
Nodes (23): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer (+15 more)

### Community 19 - "fileEditModel"
Cohesion: 0.09
Nodes (9): clipLabel(), humanSize(), paneTitle(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "newFakeVault"
Cohesion: 0.09
Nodes (57): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+49 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (14): ownerPaths(), Resource, GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles() (+6 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (42): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+34 more)

### Community 24 - "plan_test.go"
Cohesion: 0.40
Nodes (11): scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+3 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (95): add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+87 more)

### Community 26 - "highlight_test.go"
Cohesion: 0.17
Nodes (15): Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens(), TestHighlighterColorsTokens() (+7 more)

### Community 27 - "BinaryEntry"
Cohesion: 0.09
Nodes (17): BinaryEntry, diffBinaries(), TestDiffBinaries(), netGuardPill(), NetworkGuardRepository, newExeResolver(), guardRow(), clonePolicy() (+9 more)

### Community 28 - "gui/session.go"
Cohesion: 0.07
Nodes (40): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 29 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 30 - "KernelDev"
Cohesion: 0.09
Nodes (30): fdStat, BinaryStat, physicalParent(), statFD(), syncBusBrokers(), TagRelayBus(), binaryStatOf(), binaryStatOf() (+22 more)

### Community 31 - "DeployInstalledBinary"
Cohesion: 0.07
Nodes (36): installGUI(), buildBinaryIfNeeded(), ensureInstalledBinary(), installBinaryOnly(), prepareInstallation(), TestEnsureInstalledBinary(), installServices(), mustCwd() (+28 more)

### Community 32 - "GuardEvent"
Cohesion: 0.07
Nodes (25): holdHeadless(), HeadlessLine(), snapGuard, GuardEvent, logBacklogDenial(), formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts() (+17 more)

### Community 33 - "codeedit.go"
Cohesion: 0.15
Nodes (29): class, editorModel, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+21 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.18
Nodes (31): busPath(), connexec(), connFD(), createItem(), dial(), dupkey(), execInto(), failf() (+23 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (13): checkSecret(), checkTrustLine(), trustRows(), pingInterval(), File, TrustItem, backend, backendOp (+5 more)

### Community 38 - "editView"
Cohesion: 0.11
Nodes (9): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), relOf() (+1 more)

### Community 39 - "runMonitor"
Cohesion: 0.13
Nodes (18): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), pathWithinDir() (+10 more)

### Community 40 - "filevault.go"
Cohesion: 0.23
Nodes (14): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+6 more)

### Community 41 - "daemon.conf reference"
Cohesion: 0.08
Nodes (24): A complete example, `allow_env`, `allow_lib`, Binary lines and events, daemon.conf reference, `deny_args`, `warn_args`, `deny_env`, `warn_env`, `[electron_apps]` (+16 more)

### Community 42 - "FindDaemonProcesses"
Cohesion: 0.12
Nodes (18): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), anyDaemonRunning(), probeSystem() (+10 more)

### Community 44 - "pinstate.go"
Cohesion: 0.11
Nodes (33): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder() (+25 more)

### Community 45 - "fakeGuardRepo"
Cohesion: 0.07
Nodes (17): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+9 more)

### Community 46 - "updatePage"
Cohesion: 0.24
Nodes (6): installedBinary, latestRelease, updatePage, readInstalled(), boolIndex(), isPreRelease()

### Community 48 - "os.File"
Cohesion: 0.05
Nodes (59): scriptIdentity(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched() (+51 more)

### Community 49 - "fadeScroll"
Cohesion: 0.08
Nodes (16): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf(), newFadeScroll(), newTestFadeScroll(), settle() (+8 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.07
Nodes (29): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+21 more)

### Community 51 - "protectedPage"
Cohesion: 0.05
Nodes (25): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, recordReview, recordView (+17 more)

### Community 52 - "Client"
Cohesion: 0.07
Nodes (19): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+11 more)

### Community 53 - "Finder"
Cohesion: 0.15
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 54 - "session"
Cohesion: 0.15
Nodes (5): peerCred(), readFrame(), queued, senderQueue, session

### Community 55 - ".apply"
Cohesion: 0.17
Nodes (5): applier, ledgerOps, PlannedBinary, Store, vaultOps

### Community 56 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 57 - "UI"
Cohesion: 0.06
Nodes (34): nextConfig(), showChangelog(), channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI() (+26 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (33): engine, portOwner, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+25 more)

### Community 59 - "NetEventType"
Cohesion: 0.07
Nodes (18): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+10 more)

### Community 60 - "TrustGuard"
Cohesion: 0.05
Nodes (28): lsmHook, tempGrant, trustHook, trustRows, checkBtrfsKeys(), ensureBtrfsLayout(), droppableFreed(), freedSuperseded() (+20 more)

### Community 61 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 62 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 63 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.13
Nodes (21): placedWalk, treeJudge, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedEntry(), placedFileStat() (+13 more)

### Community 64 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "main_test.go"
Cohesion: 0.06
Nodes (19): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+11 more)

### Community 68 - "Open"
Cohesion: 0.08
Nodes (34): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+26 more)

### Community 69 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 70 - "Vault"
Cohesion: 0.12
Nodes (20): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+12 more)

### Community 71 - "linkView"
Cohesion: 0.20
Nodes (11): linkTarget(), newLinkView(), checkLink(), linkFixture(), TestCheckLinksAccepts(), TestCheckLinksJudgesEachTarget(), TestCheckLinksRefusesUserFiles(), TestCheckLinksSeesThePlannedTree() (+3 more)

### Community 72 - "planUpdaters"
Cohesion: 0.13
Nodes (20): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+12 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "binaryList"
Cohesion: 0.10
Nodes (14): binaryList, netGuardTab, policyEditor, policyHooks, policyTab, newBinaryList(), newPolicyTab(), verbFor() (+6 more)

### Community 75 - "Guard"
Cohesion: 0.05
Nodes (23): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, rootHandle, SectionScan, openRoot() (+15 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "vaultfs.go"
Cohesion: 0.16
Nodes (10): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, Open(), TestFileModeToUnixModeRoundTrip(), UnixModeToFileMode(), Unopened() (+2 more)

### Community 78 - "liveSession"
Cohesion: 0.17
Nodes (11): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), holdForward() (+3 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "GuardInodeKey"
Cohesion: 0.05
Nodes (31): exeRow, keptEntries, refusedReplacements, supersededBinary, tempMaskOp, tempRow, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted() (+23 more)

### Community 81 - "controlServer"
Cohesion: 0.10
Nodes (23): readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+15 more)

### Community 82 - "runGenKey"
Cohesion: 0.33
Nodes (7): confirmMasterKeyOverwrite(), runGenKey(), ensureMasterKey(), checkKeyLen(), GenerateMasterKey(), MasterKeyExists(), readOrCreateKey()

### Community 83 - "rcUser"
Cohesion: 0.10
Nodes (48): setupBunTmpdir(), listFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), rcTarget, BunTmpDir() (+40 more)

### Community 84 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "inode_dev"
Cohesion: 0.13
Nodes (22): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+14 more)

### Community 87 - "loginServer"
Cohesion: 0.08
Nodes (23): markReadyWithLogin(), hasSentinel(), lockParent(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor(), startLoginServer() (+15 more)

### Community 88 - "installPage"
Cohesion: 0.07
Nodes (17): actionGate, binSource, daemonPage, daemonState, handButton, handCheck, installAction, installPage (+9 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "buildTrustedSet"
Cohesion: 0.08
Nodes (31): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), buildTrustedSet() (+23 more)

### Community 92 - "batchInodeRows"
Cohesion: 0.18
Nodes (13): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+5 more)

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
Cohesion: 0.14
Nodes (14): procStats, statsSample, duration(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "run.go"
Cohesion: 0.36
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 101 - "sources"
Cohesion: 0.21
Nodes (5): TestGeneratedLibraryBlockRoundTrips(), lstatOf(), outsideHardlinks(), libDirDev(), sources

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "SanitizeText"
Cohesion: 0.05
Nodes (38): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted() (+30 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "ListUsers"
Cohesion: 0.09
Nodes (32): addedResources(), TestAddedResourcesSkipsAlreadyGuardedSSH(), runFacultatives(), cleanOrphanedFscrypt(), agentSocketNote(), askSSHAgentUsers(), removeFacultative(), removeFacultatives() (+24 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "Candidate"
Cohesion: 0.11
Nodes (26): appendSectionsAndEdit(), collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates() (+18 more)

### Community 110 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.09
Nodes (20): lineStyles, pos, viewportState, diffModel, reposition(), fg(), baseStyle(), gutterLabel() (+12 more)

### Community 111 - "serverAuth"
Cohesion: 0.32
Nodes (5): TestClientAuth(), TestSASL(), clientAuth(), serverAuth(), saslServer

### Community 112 - "netGuardModel"
Cohesion: 0.29
Nodes (3): protoString(), netGuardEventMsg, netGuardModel

### Community 113 - "recordFlow"
Cohesion: 0.08
Nodes (35): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), lines(), offerBackupDeletion(), pickBinaries() (+27 more)

### Community 115 - "classify"
Cohesion: 0.14
Nodes (20): call(), dictFrame(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyItemTypeIsSchema(), TestClassifyOpenCallsExact(), TestDuplicateDictKeyRefused(), TestReadFrame() (+12 more)

### Community 116 - "Recorder"
Cohesion: 0.13
Nodes (19): Target, newBinary(), drain(), Recorder, Target, monitorTargets(), newRecorder(), opNames() (+11 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "time.Time"
Cohesion: 0.25
Nodes (5): LockedError, sweepID, sweepMemo, listEntries(), dirEntry

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "User"
Cohesion: 0.10
Nodes (30): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor() (+22 more)

### Community 124 - "inet_connect.c"
Cohesion: 0.25
Nodes (17): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+9 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 127 - "EventType"
Cohesion: 0.08
Nodes (18): selfEventsFor(), snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), specFromForm() (+10 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 131 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (38): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+30 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 143 - "validateConfigText"
Cohesion: 0.11
Nodes (17): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse() (+9 more)

### Community 144 - "dialLiveSession"
Cohesion: 0.07
Nodes (29): auditAfterEdit(), auditFindings(), dialLiveSession(), applyEdit(), editAgain(), fetchConfig(), putConfig(), runEditConfig() (+21 more)

### Community 145 - "Ledger"
Cohesion: 0.09
Nodes (13): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+5 more)

### Community 146 - "MonitorRepository"
Cohesion: 0.13
Nodes (13): MonitorRepository, NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+5 more)

### Community 147 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 148 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 149 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (11): DrainUntilSignal(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+3 more)

### Community 150 - "RefuseIdentityLines"
Cohesion: 0.20
Nodes (10): loaderName(), LoaderRefusal(), soNames(), hostLoader(), TestLoaderRefusal(), RefuseIdentityLines(), TestRefuseIdentityLines(), whitelistLines() (+2 more)

### Community 152 - "watchSet"
Cohesion: 0.23
Nodes (6): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet

### Community 154 - "Monitor"
Cohesion: 0.13
Nodes (5): NewMonitor(), putBtrfsLayout(), exeKey, Monitor, watchedInode

### Community 155 - "debouncer"
Cohesion: 0.33
Nodes (4): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), debouncer

### Community 156 - "io/fs.FileInfo"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - ".admitGated"
Cohesion: 0.26
Nodes (4): placedKey(), Caller, Role, Judge

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "wire.go"
Cohesion: 0.23
Nodes (9): decodeFrame(), distinctDictKeys(), header(), sameDict(), sameEntries(), sameList(), signatureField(), skipHeaderValue() (+1 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "update_test.go"
Cohesion: 0.05
Nodes (64): installHooks(), releaseBuild(), Run(), staged(), applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions() (+56 more)

### Community 163 - "reloadOnce"
Cohesion: 0.06
Nodes (40): CheckBPFLSM(), CheckEBPF(), newControlManager(), startControlManager(), awaitStartupOrSignal(), backingDeviceUnion(), buildGuards(), catchLifecycleSignals() (+32 more)

### Community 164 - "absPath"
Cohesion: 0.07
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (48): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+40 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (52): fakeLink, pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines (+44 more)

### Community 167 - "safeio.go"
Cohesion: 0.10
Nodes (30): TestDaemonUnitWithMetadataOutput(), dupFD(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow(), TestInstallSSHAgentRefusesSymlinkedUnitDir() (+22 more)

### Community 168 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.09
Nodes (11): recordTickCmd(), recordModel, clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), ParseMode(), noticeModel (+3 more)

### Community 169 - "changelog.go"
Cohesion: 0.12
Nodes (28): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit() (+20 more)

### Community 170 - "vaultfs_test.go"
Cohesion: 0.17
Nodes (15): IsBinaryContent(), IsOctalMode(), inodeOf(), mustRead(), mustWrite(), openRootT(), TestDefaultNewFileMode(), TestIsBinaryContent() (+7 more)

### Community 171 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "cursor.go"
Cohesion: 0.07
Nodes (20): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+12 more)

### Community 174 - "PlanDaemonUnitUpdate"
Cohesion: 0.27
Nodes (10): ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate(), TestPlanDaemonUnitUpdateDropsLocalEdits(), TestShippedDaemonUnit() (+2 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.13
Nodes (15): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+7 more)

### Community 178 - "Daemon configuration"
Cohesion: 0.14
Nodes (13): `[inspectors]`, Daemon configuration, Daemon settings, Generic Electron apps, Launch rules, Library trust, Port sections, Process inspectors (+5 more)

### Community 179 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 181 - "NetGuardEvent"
Cohesion: 0.09
Nodes (12): snapNetGuard, reasonText(), NetGuardEvent, TestReasonText(), netGuardEventLine, fakeEngines, fakeGuard, fakeMonitor (+4 more)

### Community 182 - "sync.Mutex"
Cohesion: 0.11
Nodes (11): fyneLogFilter, logHook, logLine, journalLine(), TestJournalLine(), logRow(), logStyle(), formatEntry() (+3 more)

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

### Community 191 - "catalogRefresher"
Cohesion: 0.21
Nodes (5): drain(), newCatalogRefresher(), readEvents(), swapFile(), catalogRefresher

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - ".netIntegrityContainer"
Cohesion: 0.28
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 194 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 195 - "decryptStep"
Cohesion: 0.22
Nodes (9): askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), TestExistingEncryptedDirsNilPaths(), TestExistingEncryptedDirsSkipsPlain(), ExistingEncryptedDirs() (+1 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "RunWith"
Cohesion: 0.16
Nodes (15): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), stagedBinary, Mode, Request (+7 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "os.FileInfo"
Cohesion: 0.21
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 202 - "Gate"
Cohesion: 0.31
Nodes (3): Gate, Listen(), Bus

### Community 203 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 205 - "runView"
Cohesion: 0.13
Nodes (5): runHost, runView, stageItem, stageState, captureLogs()

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (33): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+25 more)

### Community 209 - "IntegrationSuite"
Cohesion: 0.23
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 210 - "RefreshSection"
Cohesion: 0.09
Nodes (30): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+22 more)

### Community 211 - "maintainer/plan.go"
Cohesion: 0.31
Nodes (7): isWhiteout(), opaque(), pathByte(), pathRecord(), TestPathRecord(), opKind, walker

### Community 212 - "UnifiedDiff"
Cohesion: 0.32
Nodes (8): checkAndApply(), refreshCatalogConfig(), installConfig(), ConfirmOverwrite(), ShowDiff(), TestUnifiedDiff(), TestUnifiedDiffIdentical(), UnifiedDiff()

### Community 213 - "watchPattern"
Cohesion: 0.52
Nodes (6): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern

### Community 214 - "sources_test.go"
Cohesion: 0.38
Nodes (4): newSources(), TestSourcesHoldFilesAndArchiveMembers(), TestSourcesRefuseOtherFilesystem(), writeTar()

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 219 - "find_test.go"
Cohesion: 0.17
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 221 - "uninstallAll"
Cohesion: 0.12
Nodes (26): deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), offerBackupCleanup(), run(), runUninstall(), RunWith() (+18 more)

### Community 222 - "ValidateTargets"
Cohesion: 0.19
Nodes (9): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets(), Existing (+1 more)

### Community 223 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 224 - "runNetworkGuard"
Cohesion: 0.10
Nodes (23): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), WithSuspectsFrom(), DiscoverInfraBinaries() (+15 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 227 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 228 - "loadDaemonConfig"
Cohesion: 0.20
Nodes (10): loadDaemonConfig(), noteReloadedSettings(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds(), SetUserHomes() (+2 more)

### Community 229 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 231 - "SecretMatch"
Cohesion: 0.40
Nodes (4): SecretScope, SecretMatch, parseSecretMatch(), TestSecretMatch()

### Community 233 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 234 - "mutedText"
Cohesion: 0.08
Nodes (39): checkGrid, filesPage, homePage, multiAnswer, panel, recordForm, statusTile, newBinaryRow() (+31 more)

### Community 235 - "filevault_test.go"
Cohesion: 0.18
Nodes (26): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestEnsureRecoverySidecarPlaceholder() (+18 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "mountTable"
Cohesion: 0.13
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), mountTable (+1 more)

### Community 238 - "GenerateConf"
Cohesion: 0.04
Nodes (69): TestDiffMergePreservesExistingAndParses(), sectionsFromCandidates(), sharedSection(), discordCatalogEntry(), groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSectionsFromCandidatesSkipsMissingWatchSubPaths() (+61 more)

### Community 241 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 247 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 252 - "netModel"
Cohesion: 0.15
Nodes (10): formatType(), padStyled(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine (+2 more)

### Community 254 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 256 - "expandPlaceholders"
Cohesion: 0.11
Nodes (11): systemPatterns(), systemRule(), expandPlaceholders(), BinaryRule, TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob (+3 more)

### Community 257 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 267 - "Config"
Cohesion: 0.04
Nodes (93): ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), parseReservedPorts(), readReservedPorts(), TestParseReservedPorts(), ensureSecretAnchors() (+85 more)

### Community 280 - ".keyringObject"
Cohesion: 0.28
Nodes (3): placeableAttributes(), TestLearnKeyringFromFirstReply(), TestPlaceableAttributes()

### Community 281 - "hashFileShared"
Cohesion: 0.21
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

## Knowledge Gaps
- **34 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck` (+29 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 688 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **49 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `fakeDaemon`, `guardUnitTest`, `.dispatch`, `buildOneGuard`, `DaemonEvent`, `hashFileShared`, `BinaryEntry`, `GuardEvent`, `reloadOnce`, `os.File`, `sync.Mutex`, `engine`, `TrustGuard`, `GuardInodeKey`, `loginServer`, `batchInodeRows`, `time.Duration`, `time.Time`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _34 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.09345187133072407 - nodes in this community are weakly interconnected._
- **Why does `shQuote()` connect `IntegrationSuite` to `go_pkg_strings`, `IntegrationSuite`, `.runEditorHarness`, `IntegrationSuite`, `absPath`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.014925373134328358 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `.netIntegrityContainer`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0943859649122807 - nodes in this community are weakly interconnected._