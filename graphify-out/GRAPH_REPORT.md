# Graph Report - app-listener  (2026-11-07)

## Corpus Check
- 575 files · ~1,821,145 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7549 nodes · 34965 edges · 255 communities (206 shown, 49 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3102 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `375c4f35`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- github.com/charmbracelet/bubbles/textarea.Model
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- .showSized
- ConfinedEntry
- RefreshSection
- Guard
- buildOneGuard
- DaemonEvent
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- op
- guard.bpf.c
- Finder
- inode_dev
- gui/session.go
- rcUser
- time.Time
- BinaryEntry
- ListUsers
- controlServer
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- backend.go
- runNetworkGuard
- desktop_test.go
- UI
- fileedit_model_test.go
- monitorUnitTest
- pinstate.go
- updatePage
- runNetworkMonitor
- stdio
- runView
- fadeScroll
- planUpdaters
- protectedPage
- Client
- highlight_test.go
- openFileVaultWithMasterKey
- .apply
- backend_test.go
- serveWebSocket
- GuardInodeKey
- IntegrationSuite
- installPage
- update/update.go
- Ledger
- golang.org/x/sys/unix.Stat_t
- daemon.conf reference
- unix_connect.c
- IntegrationSuite
- wire.go
- buildTrustedSet
- parseDaemonEvents
- Vault
- Candidate
- changelog.go
- fcntl
- FileTarget
- FileEvent
- codecheck/main.go
- NetGuard
- liveSession
- process_vm_readv.c
- .filterEvent
- safeio.go
- .ExpandLibDirWriters
- daemon_multicall_test.go
- GuardEvent
- bpfstats/main.go
- exe_supersede.h
- loginServer
- compileGlobNames
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LaunchRule
- Config
- netinject.c
- runOfflineEdit
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- recordFlow
- journal
- Well-known bypass classes
- runLockdown
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- vaultfs.go
- Build & sign release assets (reusable workflow)
- serve.go
- newTestTarget
- NewDumpHook
- audit.go
- NetEventType
- IntegrationSuite
- startCatalogRefresh
- daemon: permanent protection with encryption at rest
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- recoverFileInPlace
- IntegrationSuite
- WebSocket /ws client
- Binary
- inet_connect.c
- install.sh
- Prereq
- User
- make test-integration (rootful Docker suite)
- SelectOrphans
- finishBinary
- github.com/spf13/cobra.Command
- trace-app-libs.sh
- IntegrationSuite
- github.com/charmbracelet/bubbletea.Cmd
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- session
- KernelDev
- Monitor
- hashCache
- TestMonitorUseCaseLifecycle
- maintainerKeeper
- btrfsLayoutFrom
- Recorder
- dialLiveSession
- fakeStore
- tempGrantSetup
- IntegrationSuite
- os.FileInfo
- go_pkg_strings
- CopyTreeWithProgress
- run.go
- Installation
- NetGuardEvent
- graphify-refresh.sh
- .Layout
- backups.go
- batchInodeRows
- .populateWatchInodes
- sync.Mutex
- shell
- fakeDaemon
- github.com/stretchr/testify/suite.Suite
- PlanDaemonUnitUpdate
- LibraryClosure
- writeEvent
- Limitations
- eventUnitTest
- plan_test.go
- net_tester/main.go
- editauth_test.go
- Plan
- openBunTmp
- IntegrationSuite
- New
- TrustGuard
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- IntegrationSuite
- os.File
- Self-updating apps and the catalog
- serveHandler
- WithBottomBar
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- IntegrationSuite
- pathCache
- loadDaemonConfig
- refuseOutsideHardlinks
- .netIntegrityContainer
- testLimiter
- EventType
- monitor
- IntegrationSuite
- absPath
- linkView
- runMonitor
- catalogRefresher
- update_test.go
- FindDaemonProcesses
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- NewEventFanout
- RootPlacedName
- find_test.go
- VouchedDevs
- DiscoverInfraBinaries
- IntegrationSuite
- inodeChain
- Plan
- supersede_probe.c
- ShortRefused
- SecretMatch
- launch_probe.c
- IntegrationSuite
- github.com/Virgula0/app-listener/tools/codecheck
- procstats.go
- taint_victim.c
- parseResize
- mutedText
- codeedit.go
- Troubleshooting
- mountTable
- GenerateConf
- CleanupStalePins
- IntegrationSuite
- Vault
- .runEditorHarness
- daemon-config.md
- fetchReleases
- Development
- app-listener documentation
- LiveEmptyWhitelistRejected
- go_pkg_path_filepath
- monitor: see what touches a path
- network monitor: see what a program talks to
- .eventTypeAllowed
- bleedLayout
- daemonconfig.go

## God Nodes (most connected - your core abstractions)
1. `Config` - 149 edges
2. `Guard` - 126 edges
3. `GuardInodeKey` - 116 edges
4. `Load()` - 101 edges
5. `absPath()` - 97 edges
6. `shell` - 97 edges
7. `IntegrationSuite` - 83 edges
8. `writeConfig()` - 79 edges
9. `fileEditModel` - 77 edges
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

## Communities (255 total, 49 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.16
Nodes (29): tempJournalRow, argvMatch, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion (+21 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (161): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload() (+153 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote(), watchPathsInConfig()

### Community 3 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (17): lineStyles, pos, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+9 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (85): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+77 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.23
Nodes (5): guardBinaryFlag(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (24): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+16 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (38): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetHomeDepartureTriggersRefresh(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+30 more)

### Community 10 - ".showSized"
Cohesion: 0.14
Nodes (7): installUI, resizeWatch, sizedDialog, ask(), diffText(), promptBody(), showAudit()

### Community 11 - "ConfinedEntry"
Cohesion: 0.12
Nodes (28): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+20 more)

### Community 12 - "RefreshSection"
Cohesion: 0.07
Nodes (35): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+27 more)

### Community 13 - "Guard"
Cohesion: 0.04
Nodes (26): backingDeviceUnion(), deferredBinary, exeRow, rootHandle, supersededBinary, tempMaskOp, tempRow, TestAddBinaryActionsRefusesBeforeWriting() (+18 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.05
Nodes (63): DrainUntilSignal(), CheckEBPF(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration(), prepareDaemonStart(), reloadOnce() (+55 more)

### Community 15 - "DaemonEvent"
Cohesion: 0.18
Nodes (17): ServeConfig, catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+9 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.08
Nodes (69): argv_env(), argv_matched(), argv_not_root(), argv_rule_at(), argv_short(), argv_word(), caller_updates(), creds_overridden() (+61 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (22): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer (+14 more)

### Community 19 - "fileEditModel"
Cohesion: 0.07
Nodes (13): newFindInput(), clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser (+5 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.13
Nodes (26): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+18 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.13
Nodes (52): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+44 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (12): Resource, narrowInLockedVault(), GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, TestRepeatLogOncePerErrorAndOnRecovery(), uniqueEncryptionRoots() (+4 more)

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (37): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+29 more)

### Community 24 - "op"
Cohesion: 0.19
Nodes (11): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), removeAll() (+3 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (97): exe_refused(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+89 more)

### Community 26 - "Finder"
Cohesion: 0.17
Nodes (8): findButton, span, findAll(), Finder, lowerRunes(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 27 - "inode_dev"
Cohesion: 0.43
Nodes (4): ctx_ptr(), inode_dev(), sb_dev(), guard_inode_free()

### Community 28 - "gui/session.go"
Cohesion: 0.07
Nodes (40): invokingIDs(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 29 - "rcUser"
Cohesion: 0.10
Nodes (47): listFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), rcTarget, BunTmpDir(), TestHasRCBlocksTrackEnsureAndRemove() (+39 more)

### Community 30 - "time.Time"
Cohesion: 0.08
Nodes (13): LockedError, sweepID, sweepMemo, codeEntry, editView, handButton, uninstallCard, Entry (+5 more)

### Community 31 - "BinaryEntry"
Cohesion: 0.08
Nodes (16): keptEntries, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath(), diffBinaries(), TestDiffBinaries(), netGuardPill() (+8 more)

### Community 32 - "ListUsers"
Cohesion: 0.11
Nodes (25): installFacultatives(), installSSHAgentFacultative(), runFacultatives(), removeFacultative(), removeFacultatives(), TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(), detectSSHAgentUnits(), removeBunLauncherEnv() (+17 more)

### Community 33 - "controlServer"
Cohesion: 0.08
Nodes (27): readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+19 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (31): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+23 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.18
Nodes (31): busPath(), connexec(), connFD(), createItem(), dial(), dupkey(), execInto(), failf() (+23 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (14): entryOf(), checkTrustLine(), File, TrustItem, backend, backendOp, backendState, TestCleanRel() (+6 more)

### Community 38 - "backend.go"
Cohesion: 0.04
Nodes (90): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), TestPickEncryptedDirSingleSkipsPicker(), pickDir(), pickEncryptedDir(), offerDaemonStart(), runDiffCatalog() (+82 more)

### Community 39 - "runNetworkGuard"
Cohesion: 0.15
Nodes (16): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+8 more)

### Community 40 - "desktop_test.go"
Cohesion: 0.07
Nodes (25): TestEveryInstallModeMapped(), buttonCursor, logGrid, multiAnswer, pageDef, stageItem, stageState, textAnswer (+17 more)

### Community 41 - "UI"
Cohesion: 0.09
Nodes (19): showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), MultiSelectKeymap() (+11 more)

### Community 42 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 44 - "pinstate.go"
Cohesion: 0.19
Nodes (21): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+13 more)

### Community 45 - "updatePage"
Cohesion: 0.10
Nodes (10): handCheckGroup, handRadioGroup, handSelect, latestRelease, updatePage, handCursor(), newNetEventChecks(), boolIndex() (+2 more)

### Community 46 - "runNetworkMonitor"
Cohesion: 0.11
Nodes (11): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runNetworkMonitor(), runTUI(), NetworkMonitorRepository, NewNetModel() (+3 more)

### Community 48 - "runView"
Cohesion: 0.08
Nodes (12): logHook, logLine, runView, journalLine(), TestJournalLine(), asset(), coloredIcon(), logRow() (+4 more)

### Community 49 - "fadeScroll"
Cohesion: 0.08
Nodes (16): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf(), newFadeScroll(), newTestFadeScroll(), settle() (+8 more)

### Community 50 - "planUpdaters"
Cohesion: 0.13
Nodes (20): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+12 more)

### Community 51 - "protectedPage"
Cohesion: 0.05
Nodes (26): backendLink, checkGrid, configView, forwardForm, forwardView, handCheck, protectedPage, recordReview (+18 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (15): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+7 more)

### Community 53 - "highlight_test.go"
Cohesion: 0.16
Nodes (17): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+9 more)

### Community 54 - "openFileVaultWithMasterKey"
Cohesion: 0.17
Nodes (15): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), openFileVaultWithMasterKey(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated() (+7 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (15): Apply(), DeleteBackup(), Outcome, ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+7 more)

### Community 56 - "backend_test.go"
Cohesion: 0.25
Nodes (21): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+13 more)

### Community 57 - "serveWebSocket"
Cohesion: 0.21
Nodes (9): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop() (+1 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (62): engine, portOwner, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), confineConnect() (+54 more)

### Community 59 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 60 - "installPage"
Cohesion: 0.10
Nodes (14): actionGate, binSource, installAction, InstallMode, installPage, InstallRequest, Options, UninstallHooks (+6 more)

### Community 61 - "update/update.go"
Cohesion: 0.13
Nodes (38): applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI() (+30 more)

### Community 62 - "Ledger"
Cohesion: 0.07
Nodes (30): trustRows(), candidates(), confirmAll(), describe(), record(), run(), `allow_env`, `deny_args`, `warn_args` (+22 more)

### Community 63 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.16
Nodes (16): placedWalk, treeJudge, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedEntry(), placedFileStat() (+8 more)

### Community 64 - "daemon.conf reference"
Cohesion: 0.06
Nodes (35): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+27 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "IntegrationSuite"
Cohesion: 0.13
Nodes (4): exploitTest, IntegrationSuite, newEventTypes(), IntegrationSuite

### Community 67 - "wire.go"
Cohesion: 0.09
Nodes (30): call(), dictFrame(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyItemTypeIsSchema(), TestClassifyOpenCallsExact(), TestDuplicateDictKeyRefused(), TestLearnKeyringFromFirstReply() (+22 more)

### Community 68 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 69 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 70 - "Vault"
Cohesion: 0.13
Nodes (18): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+10 more)

### Community 71 - "Candidate"
Cohesion: 0.12
Nodes (23): collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), groupCandidates(), pickDirectories(), pickFromCandidates(), TestGroupCandidates(), Candidate (+15 more)

### Community 72 - "changelog.go"
Cohesion: 0.17
Nodes (22): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit() (+14 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "FileTarget"
Cohesion: 0.06
Nodes (26): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+18 more)

### Community 75 - "FileEvent"
Cohesion: 0.12
Nodes (6): snapMonitor, FileEvent, MonitorRepository, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "NetGuard"
Cohesion: 0.08
Nodes (13): lsmHook, foreignPinnedLink(), PumpRingbuf(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs() (+5 more)

### Community 78 - "liveSession"
Cohesion: 0.13
Nodes (15): checkSecret(), TestCheckSecretRefusesLineBreaks(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel() (+7 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 81 - "safeio.go"
Cohesion: 0.06
Nodes (47): setupBunTmpdir(), deploy(), deployFiles(), grantNewBinaries(), preflightDeployedBinary(), TestDaemonUnitWithMetadataOutput(), dupFD(), installConfig() (+39 more)

### Community 82 - ".ExpandLibDirWriters"
Cohesion: 0.14
Nodes (13): homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelistSymlinkEscape(), TestSteamLibDirsCoverProtonCodeOnly(), TestSteamLibDirsExcludeScoutRuntime(), TestSteamLibDirWritersCoverItsOwnBinaries(), boundedGlob() (+5 more)

### Community 83 - "daemon_multicall_test.go"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 84 - "GuardEvent"
Cohesion: 0.04
Nodes (31): readTempJournal(), holdHeadless(), HeadlessLine(), forwardRequest, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents() (+23 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.13
Nodes (23): CheckBPFLSM(), runBPFCheck(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity() (+15 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.22
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 87 - "loginServer"
Cohesion: 0.09
Nodes (20): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), lockParent(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel() (+12 more)

### Community 88 - "compileGlobNames"
Cohesion: 0.29
Nodes (5): GlobChild, compileGlobNames(), globText(), resolveChildren(), TestCompileGlobNames_RepeatedPatternOrsBits()

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LaunchRule"
Cohesion: 0.14
Nodes (20): ShortOpts, argvEntry, LaunchRule, EncodeLaunchRules(), canonLaunchRules(), encodeArgvRule(), encodeEnvRule(), encodeLaunchSet() (+12 more)

### Community 92 - "Config"
Cohesion: 0.06
Nodes (55): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), askBunTmpdirUsers() (+47 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "runOfflineEdit"
Cohesion: 0.11
Nodes (18): runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), CleanOrphanedMetadata() (+10 more)

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
Cohesion: 0.13
Nodes (15): procStats, statsSample, duration(), formatClock(), cpuTime(), humanBytes(), rssBytes(), New() (+7 more)

### Community 100 - "recordFlow"
Cohesion: 0.16
Nodes (12): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+4 more)

### Community 101 - "journal"
Cohesion: 0.14
Nodes (7): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "runLockdown"
Cohesion: 0.14
Nodes (16): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), dirIsRootOwnedSafe(), isBpffs() (+8 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "vaultfs.go"
Cohesion: 0.09
Nodes (24): DefaultNewFileMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode(), inodeOf(), mustRead() (+16 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "serve.go"
Cohesion: 0.27
Nodes (10): isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), Serve(), serveHTTP(), waitProgram() (+2 more)

### Community 110 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 111 - "NewDumpHook"
Cohesion: 0.19
Nodes (15): TestSetupDumpLogRefusesSymlink(), applyVerbosity(), effectiveVerbose(), preRun(), setupDumpLog(), VerboseLevel, BridgeStdLog(), NewDumpHook() (+7 more)

### Community 112 - "audit.go"
Cohesion: 0.16
Nodes (12): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+4 more)

### Community 113 - "NetEventType"
Cohesion: 0.08
Nodes (15): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+7 more)

### Community 115 - "startCatalogRefresh"
Cohesion: 0.21
Nodes (6): newControlManager(), startControlManager(), makeReloadHandler(), startCatalogRefresh(), controlManager, SetSupersedeHold()

### Community 116 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.13
Nodes (15): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+7 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "recoverFileInPlace"
Cohesion: 0.31
Nodes (13): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), recoverFileInPlace(), rejectSidecarSymlink(), stageRecovery(), inode() (+5 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "Binary"
Cohesion: 0.18
Nodes (13): formatElapsed(), lines(), orNone(), pickBinaries(), pickLibs(), printBinary(), printRecording(), reviewRecording() (+5 more)

### Community 124 - "inet_connect.c"
Cohesion: 0.25
Nodes (17): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+9 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 127 - "User"
Cohesion: 0.07
Nodes (32): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor() (+24 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 130 - "finishBinary"
Cohesion: 0.19
Nodes (12): beyondEvents(), finishBinary(), isShellHelper(), newBinary(), noteWhitelisted(), opNames(), sortedKeys(), TestFileReadersStartDeselected() (+4 more)

### Community 131 - "github.com/spf13/cobra.Command"
Cohesion: 0.19
Nodes (16): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+8 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (48): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), formatGuardEventLine(), sanitizeTerminalText() (+40 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "session"
Cohesion: 0.05
Nodes (20): placedKey(), secretJudge, SecretService, LiveGuardAt(), connectionGone(), Caller, Gate, Role (+12 more)

### Community 143 - "KernelDev"
Cohesion: 0.06
Nodes (41): BusClient, fdStat, rootSlot, openRoot(), physicalParent(), slotOf(), statAt(), statFD() (+33 more)

### Community 144 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 145 - "hashCache"
Cohesion: 0.12
Nodes (17): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, HashCache (+9 more)

### Community 146 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 147 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 148 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 149 - "Recorder"
Cohesion: 0.17
Nodes (14): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+6 more)

### Community 150 - "dialLiveSession"
Cohesion: 0.08
Nodes (32): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply() (+24 more)

### Community 152 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 154 - "os.FileInfo"
Cohesion: 0.21
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 155 - "go_pkg_strings"
Cohesion: 0.05
Nodes (49): installHooks(), releaseBuild(), Run(), staged(), atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent() (+41 more)

### Community 156 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "run.go"
Cohesion: 0.32
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "NetGuardEvent"
Cohesion: 0.18
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "backups.go"
Cohesion: 0.15
Nodes (24): deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), offerBackupCleanup(), run(), runUninstall(), RunWith() (+16 more)

### Community 163 - "batchInodeRows"
Cohesion: 0.27
Nodes (9): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+1 more)

### Community 165 - "sync.Mutex"
Cohesion: 0.05
Nodes (47): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+39 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (56): daemonPage, daemonState, fakeLink, pageID, reel, runHost, shell, snapEngines (+48 more)

### Community 169 - "PlanDaemonUnitUpdate"
Cohesion: 0.27
Nodes (10): ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate(), TestPlanDaemonUnitUpdateDropsLocalEdits(), TestShippedDaemonUnit() (+2 more)

### Community 170 - "LibraryClosure"
Cohesion: 0.07
Nodes (29): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+21 more)

### Community 171 - "writeEvent"
Cohesion: 0.33
Nodes (8): UIDResolver, NewUIDResolver(), writeEvent(), allowedEvent(), deniedEvent(), TestWriteEventAllowed(), TestWriteEventAllowedSuppressed(), TestWriteEventBlocked()

### Community 172 - "Limitations"
Cohesion: 0.17
Nodes (12): Already-running processes, Binaries you can write, Commands that stop working, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root (+4 more)

### Community 173 - "eventUnitTest"
Cohesion: 0.27
Nodes (4): BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent()

### Community 174 - "plan_test.go"
Cohesion: 0.09
Nodes (34): newLinkView(), checkLink(), linkFixture(), TestCheckLinksAccepts(), TestCheckLinksJudgesEachTarget(), TestCheckLinksRefusesUserFiles(), TestCheckLinksSeesThePlannedTree(), isWhiteout() (+26 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.19
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "editauth_test.go"
Cohesion: 0.11
Nodes (26): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+18 more)

### Community 178 - "Plan"
Cohesion: 0.20
Nodes (6): dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule

### Community 179 - "openBunTmp"
Cohesion: 0.33
Nodes (9): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp (+1 more)

### Community 181 - "New"
Cohesion: 0.18
Nodes (18): askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse(), harnessSuite, looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails() (+10 more)

### Community 182 - "TrustGuard"
Cohesion: 0.05
Nodes (28): setConfinedHomes(), argvMatches, statFunc, tempGrant, trustHook, trustRows, singleBinaryPaths(), setSupersedeTrust() (+20 more)

### Community 186 - "Binary trust"
Cohesion: 0.20
Nodes (10): Binary ledger, Binary replacement, Binary trust, Browser profiles, Hash cache, Launch rules, Launch scan, Multicall binaries (+2 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 191 - "os.File"
Cohesion: 0.03
Nodes (80): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome() (+72 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "serveHandler"
Cohesion: 0.31
Nodes (8): basicAuth(), securityHeaders(), serveHandler(), TestBasicAuth(), TestSecurityHeadersAndExactRoot(), TestServeHandlerProtectsPageAndWebSocket(), TestServeWebSocketSnapshotsResizeAndSingleViewer(), validateHost()

### Community 194 - "WithBottomBar"
Cohesion: 0.31
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "loadDaemonConfig"
Cohesion: 0.14
Nodes (16): loadDaemonConfig(), noteReloadedSettings(), ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), parseReservedPorts(), readReservedPorts() (+8 more)

### Community 202 - "refuseOutsideHardlinks"
Cohesion: 0.24
Nodes (8): lstatOf(), outsideHardlinks(), refuseOutsideHardlinks(), TestRefuseOutsideHardlinks_FileRoot(), TestRefuseOutsideHardlinks_InTreeLinksAccepted(), TestRefuseOutsideHardlinks_MissingRootFails(), TestRefuseOutsideHardlinks_OutsideLinkRefused(), writeFile()

### Community 204 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 205 - "EventType"
Cohesion: 0.10
Nodes (18): selfEventsFor(), bpfGuardEvent, EventType, parseEvents(), decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel() (+10 more)

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (33): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+25 more)

### Community 209 - "absPath"
Cohesion: 0.05
Nodes (12): IntegrationSuite, IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath() (+4 more)

### Community 210 - "linkView"
Cohesion: 0.43
Nodes (3): linkTarget(), entryKind, linkView

### Community 211 - "runMonitor"
Cohesion: 0.19
Nodes (14): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+6 more)

### Community 212 - "catalogRefresher"
Cohesion: 0.14
Nodes (9): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), catalogRefresher (+1 more)

### Community 213 - "update_test.go"
Cohesion: 0.08
Nodes (24): newerThanInstalled(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchGUIWithoutAsset(), TestFetchReleaseByTag(), TestFilterChannel() (+16 more)

### Community 214 - "FindDaemonProcesses"
Cohesion: 0.15
Nodes (16): anyDaemonRunning(), DaemonRunning(), FindDaemonProcesses(), isDaemonCmdline(), systemctlActive(), TestFindDaemonProcessesEmpty(), TestIsDaemonCmdline(), InstallCert (+8 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 218 - "RootPlacedName"
Cohesion: 0.29
Nodes (6): PathUnder(), RootPlacedName(), rootOnlyTree(), TestRootPlacedName(), userAssignable(), userMountable()

### Community 219 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 220 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 221 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 223 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 224 - "Plan"
Cohesion: 0.38
Nodes (3): openDir(), libPlan, Plan

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "ShortRefused"
Cohesion: 0.33
Nodes (5): ShortRefused(), TestSSHConfigFileRefusedInEveryCluster(), shortMask(), shortHit(), TestShortMaskRefusesNonLetters()

### Community 227 - "SecretMatch"
Cohesion: 0.40
Nodes (4): SecretScope, SecretMatch, parseSecretMatch(), TestSecretMatch()

### Community 228 - "launch_probe.c"
Cohesion: 0.46
Nodes (6): anchor(), copy(), is_verb(), main(), read_once(), usage()

### Community 231 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 233 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 234 - "mutedText"
Cohesion: 0.06
Nodes (44): binaryList, filesPage, homePage, panel, recordForm, newBinaryList(), newBinaryRow(), newCheck() (+36 more)

### Community 235 - "codeedit.go"
Cohesion: 0.15
Nodes (31): nextConfig(), class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+23 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "mountTable"
Cohesion: 0.22
Nodes (7): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), ReadMountinfo(), unescapeMountinfo(), mountTable

### Community 238 - "GenerateConf"
Cohesion: 0.04
Nodes (87): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), editConfig(), libraryBlocksFromCandidates(), runConfigEditor(), sectionsFromCandidates(), sharedSection() (+79 more)

### Community 239 - "CleanupStalePins"
Cohesion: 0.40
Nodes (5): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 244 - "fetchReleases"
Cohesion: 0.22
Nodes (8): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError(), progressReader

### Community 245 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 246 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 255 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 256 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 257 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 267 - "daemonconfig.go"
Cohesion: 0.06
Nodes (77): systemPatterns(), systemRule(), runtimePathsFor(), watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBinaryLine() (+69 more)

## Knowledge Gaps
- **34 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck` (+29 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 691 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **49 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `time.Duration`, `batchInodeRows`, `sync.Mutex`, `fakeDaemon`, `guardUnitTest`, `ConfinedEntry`, `daemonconfig.go`, `buildOneGuard`, `DaemonEvent`, `session`, `KernelDev`, `GuardEvent`, `TrustGuard`, `GuardInodeKey`, `BinaryEntry`, `time.Time`, `os.File`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _34 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01720386236910105 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `IntegrationSuite`, `.netIntegrityContainer`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `go_pkg_strings`, `IntegrationSuite`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0943859649122807 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `Installation`, `app-listener documentation`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Should `github.com/charmbracelet/bubbles/textarea.Model` be split into smaller, more focused modules?**
  _Cohesion score 0.1026827012025902 - nodes in this community are weakly interconnected._