# Graph Report - app-listener  (2026-10-28)

## Corpus Check
- 526 files · ~1,745,254 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6834 nodes · 30274 edges · 245 communities (201 shown, 44 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2781 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e746187d`
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
- TempBinary
- New
- protectedPage
- os.File
- go_pkg_path_filepath
- installPage
- IntegrationSuite
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- Guard
- guard.bpf.c
- sync.Mutex
- Finder
- Network
- golang.org/x/sys/unix.Stat_t
- github.com/charmbracelet/bubbletea.Cmd
- TestMonitorUseCaseLifecycle
- policyEditor
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
- unit.go
- fileedit_model_test.go
- journal
- stdio
- uninstall_system.go
- fadeScroll
- vaultfs.go
- updatePage
- Client
- daemonconfig.go
- Vault
- .apply
- .FilterExistingWhitelist
- mutedText
- GuardInodeKey
- runNetworkGuard
- os.FileInfo
- Recorder
- exe_supersede.h
- compileGlobNames
- Serve
- github.com/charmbracelet/bubbles/textarea.Model
- wm.go
- IntegrationSuite
- User
- liveSession
- codeEntry
- RefreshSection
- FileEvent
- string
- bpfstats/main.go
- CandidateDir
- codecheck/main.go
- rcUser
- update_test.go
- process_vm_readv.c
- GuardEvent
- maintainer/plan.go
- Candidate
- logLine
- gui/session.go
- Monitor
- planUpdaters
- dialLiveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- EventType
- newTestTarget
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- BinaryEntry
- pinstate.go
- Well-known bypass classes
- NetGuardEvent
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- .populateWatchInodes
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- NetworkMonitor
- github.com/stretchr/testify/suite.Suite
- githubRelease
- recordFlow
- IntegrationSuite
- procstats.go
- newTestVetter
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
- root.go
- LibraryClosure
- trace-app-libs.sh
- IntegrationSuite
- maintainerKeeper
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- catalogRefresher
- serveWebSocket
- github.com/spf13/cobra.Command
- .showSized
- Ledger
- KernelDev
- Daemon configuration
- changelog.go
- Config
- fakeStore
- IntegrationSuite
- IntegrationSuite
- buildOneGuard
- safeio.go
- deploy
- daemonPage
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- dialog_test.go
- IntegrationSuite
- .netIntegrityContainer
- feed.go
- shell
- testKeyPair
- time.Time
- io.Writer
- NewEventFanout
- runLockdown
- Limitations
- runBackend
- IntegrationSuite
- runMonitor
- netinject.c
- CopyTreeWithProgress
- RawTarget
- main_test.go
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
- finishBinary
- run.go
- edit-protected
- plan_test.go
- kernel-prof.sh
- FetchLatest
- pathCache
- appTheme
- DaemonEvent
- Use
- mountTable
- uninstall/uninstall.go
- Desktop GUI
- IntegrationSuite
- guardModel
- loadDaemonConfig
- .eventTypeAllowed
- testLimiter
- inode_dev
- prctl_gate
- .awaitDaemonUpOrDead
- Vault
- watchPattern
- tempGrantSetup
- run
- shQuote
- NetEventType
- Active
- selfGuards
- install/install.go
- supersede_probe.c
- securewrite_test.go
- inspector_probe.c
- openBunTmp
- links_test.go
- github.com/Virgula0/app-listener/tools/codecheck
- Prereq
- Existing
- taint_victim.c
- WithBottomBar
- logGrid
- Troubleshooting
- linkView
- .exeKeys
- VouchedDevs
- bleedLayout
- parsePasswd
- TestCleanRel
- TestCatalogSanity
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
- `certifyInstall()` --calls--> `InstallCertPath()`  [EXTRACTED]
  cmd/functions/install/install.go → internal/protected/installcert.go

## Import Cycles
- None detected.

## Communities (245 total, 44 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.12
Nodes (36): atomicWriteAt(), writeAndSync(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs (+28 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (143): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), selfProtectSpecs(), TestAwaitStartupOrSignalAbortsAndLocksDown() (+135 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "wipe"
Cohesion: 0.13
Nodes (17): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn() (+9 more)

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
Nodes (20): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+12 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.15
Nodes (37): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+29 more)

### Community 10 - "TempBinary"
Cohesion: 0.07
Nodes (31): BinaryStat, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow() (+23 more)

### Community 11 - "New"
Cohesion: 0.09
Nodes (46): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+38 more)

### Community 12 - "protectedPage"
Cohesion: 0.04
Nodes (34): backendLink, configView, editView, forwardForm, forwardView, protectedPage, recordForm, recordReview (+26 more)

### Community 13 - "os.File"
Cohesion: 0.04
Nodes (76): DrainUntilSignal(), multicallAdmissible(), ownerPaths(), scriptIdentity(), drainEphemeral(), runHeadless(), vetDecision, VettedInode (+68 more)

### Community 15 - "installPage"
Cohesion: 0.06
Nodes (26): actionGate, binSource, installAction, installedBinary, InstallMode, installPage, InstallRequest, Options (+18 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "runView"
Cohesion: 0.03
Nodes (27): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+19 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (13): clipLabel(), humanSize(), listEntries(), paneTitle(), sortNodes(), screen(), screenCell(), chownUser (+5 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.16
Nodes (46): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+38 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (12): fakeDaemon, Resource, GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles() (+4 more)

### Community 23 - "auth.go"
Cohesion: 0.11
Nodes (42): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+34 more)

### Community 24 - "Guard"
Cohesion: 0.06
Nodes (19): deferredBinary, exeRow, rootHandle, supersededBinary, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), SharedPinDegraded(), TestNoteRootMountFollowsMountID() (+11 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (59): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with() (+51 more)

### Community 26 - "sync.Mutex"
Cohesion: 0.06
Nodes (18): lsmHook, refusedReplacements, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, setSupersedeTrust() (+10 more)

### Community 27 - "Finder"
Cohesion: 0.15
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 28 - "Network"
Cohesion: 0.06
Nodes (21): Target, guardPill(), NetworkGuardRepository, newExeResolver(), Files, FileSpec, FileTarget, Policy (+13 more)

### Community 29 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.09
Nodes (28): placedWalk, treeJudge, FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs() (+20 more)

### Community 30 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (16): recordTickCmd(), recordModel, diffModel, editorModel, clamp(), listenForGuardEvents(), tickGuardStats(), renderBullets() (+8 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.19
Nodes (11): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+3 more)

### Community 32 - "policyEditor"
Cohesion: 0.14
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+4 more)

### Community 33 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 34 - "picker"
Cohesion: 0.10
Nodes (18): keyEntry, listOpts, pickEntry, picker, treeSearch, errText(), isFindShortcut(), newKeyEntry() (+10 more)

### Community 35 - "fscrypt/orphans.go"
Cohesion: 0.19
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (12): checkSecret(), TestLiveSessionWatchAndPinger(), holdForward(), pingInterval(), File, backend, backendOp, backendState (+4 more)

### Community 38 - "controlServer"
Cohesion: 0.05
Nodes (42): swapFile(), readConfigPut(), newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload() (+34 more)

### Community 39 - "codeedit.go"
Cohesion: 0.10
Nodes (43): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+35 more)

### Community 40 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "unit.go"
Cohesion: 0.06
Nodes (49): checkAndApply(), runDiffCatalog(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary() (+41 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.23
Nodes (23): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+15 more)

### Community 46 - "journal"
Cohesion: 0.14
Nodes (8): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 48 - "uninstall_system.go"
Cohesion: 0.09
Nodes (16): detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertBunLaunchers(), revertSSHAgents() (+8 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+19 more)

### Community 50 - "vaultfs.go"
Cohesion: 0.08
Nodes (29): checkFreshDir(), openDirNoFollow(), renewDir(), entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent() (+21 more)

### Community 51 - "updatePage"
Cohesion: 0.09
Nodes (16): buttonCursor, handCheckGroup, handRadioGroup, handSelect, latestRelease, updatePage, handCursor(), newButtonCursor() (+8 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (13): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+5 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.13
Nodes (40): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+32 more)

### Community 54 - "Vault"
Cohesion: 0.07
Nodes (33): newBinaryVetter(), openBinaryVetter(), backingDeviceUnion(), buildGuards(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon(), startGuardedDaemonAbortable() (+25 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (15): Apply(), DeleteBackup(), Outcome, ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+7 more)

### Community 56 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (40): binaryList, checkGrid, filesPage, handButton, handCheck, uninstallCard, newBinaryList(), newBinaryRow() (+32 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (71): markReady(), notifySystemdReady(), engine, inodeBatchFn, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank() (+63 more)

### Community 59 - "runNetworkGuard"
Cohesion: 0.15
Nodes (16): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+8 more)

### Community 60 - "os.FileInfo"
Cohesion: 0.19
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 61 - "Recorder"
Cohesion: 0.14
Nodes (16): drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), record() (+8 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.20
Nodes (16): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+8 more)

### Community 63 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 64 - "Serve"
Cohesion: 0.17
Nodes (17): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+9 more)

### Community 65 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (19): lineStyles, pos, viewportState, reposition(), Highlighter, baseStyle(), gutterLabel(), lineLen() (+11 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 68 - "User"
Cohesion: 0.11
Nodes (35): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), applyDiffAdditions(), bunLaunchersOffered(), installFacultatives() (+27 more)

### Community 69 - "liveSession"
Cohesion: 0.22
Nodes (10): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck(), normalizedEvents() (+2 more)

### Community 71 - "RefreshSection"
Cohesion: 0.10
Nodes (26): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+18 more)

### Community 72 - "FileEvent"
Cohesion: 0.12
Nodes (6): snapMonitor, FileEvent, MonitorRepository, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 73 - "string"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "CandidateDir"
Cohesion: 0.07
Nodes (29): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), Discover() (+21 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "rcUser"
Cohesion: 0.09
Nodes (52): setupBunTmpdir(), listFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), removeFacultative(), removeFacultatives() (+44 more)

### Community 78 - "update_test.go"
Cohesion: 0.09
Nodes (15): downloadFile(), TestAssetsFor(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchGUIWithoutAsset(), TestFilterChannel(), TestGUIAssetsFor(), TestNewerThanStable() (+7 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.05
Nodes (20): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), TestLocalizePaths() (+12 more)

### Community 81 - "maintainer/plan.go"
Cohesion: 0.14
Nodes (17): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+9 more)

### Community 82 - "Candidate"
Cohesion: 0.31
Nodes (11): appendSectionsAndEdit(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickFromCandidates(), runConfigEditor(), sectionsFromCandidates() (+3 more)

### Community 83 - "logLine"
Cohesion: 0.14
Nodes (10): logHook, logLine, journalLine(), TestJournalLine(), logRow(), logStyle(), formatEntry(), sortedKeys() (+2 more)

### Community 84 - "gui/session.go"
Cohesion: 0.20
Nodes (15): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+7 more)

### Community 86 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 87 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 88 - "dialLiveSession"
Cohesion: 0.06
Nodes (37): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditTree(), coveredByWatch(), relOrPath(), TestAuditTree() (+29 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "EventType"
Cohesion: 0.10
Nodes (19): BpfEvent, EventType, eventUnitTest, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock(), specFromForm() (+11 more)

### Community 92 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (11): PumpRingbuf(), reasonText(), eventTypeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), TestReasonText() (+3 more)

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
Cohesion: 0.08
Nodes (39): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath() (+31 more)

### Community 101 - "pinstate.go"
Cohesion: 0.12
Nodes (30): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState() (+22 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NetGuardEvent"
Cohesion: 0.17
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

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

### Community 110 - "NetworkMonitor"
Cohesion: 0.14
Nodes (7): BinariesSummary(), TestBinariesSummary(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 111 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.16
Nodes (4): harnessSuite, mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 112 - "githubRelease"
Cohesion: 0.23
Nodes (17): newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime(), newerThanInstalled(), newerThanStable() (+9 more)

### Community 113 - "recordFlow"
Cohesion: 0.12
Nodes (22): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), lines(), offerBackupDeletion(), pickBinaries() (+14 more)

### Community 115 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 116 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (61): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+53 more)

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

### Community 130 - "root.go"
Cohesion: 0.07
Nodes (25): releaseBuild(), Run(), closeDumpFile(), postRun(), Progress, User, invalidPassword, stageItem (+17 more)

### Community 131 - "LibraryClosure"
Cohesion: 0.07
Nodes (28): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+20 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "maintainerKeeper"
Cohesion: 0.20
Nodes (10): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp, maintainerKeeper (+2 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "catalogRefresher"
Cohesion: 0.17
Nodes (7): drain(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), catalogRefresher, debouncer

### Community 143 - "serveWebSocket"
Cohesion: 0.15
Nodes (12): configureWebSocket(), decodeEndsAtEOF(), parseResize(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestParseResize(), TestSameOrigin() (+4 more)

### Community 144 - "github.com/spf13/cobra.Command"
Cohesion: 0.09
Nodes (34): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+26 more)

### Community 145 - ".showSized"
Cohesion: 0.10
Nodes (14): installUI, multiAnswer, resizeWatch, sizedDialog, scrollable(), selectable(), fadeScrollOf(), newFadeScroll() (+6 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "KernelDev"
Cohesion: 0.18
Nodes (14): KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), codeVouched(), isTraced(), loaderEnv(), readMountDevs() (+6 more)

### Community 148 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 149 - "changelog.go"
Cohesion: 0.22
Nodes (12): changelogText(), cumulativeChangelog(), newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize() (+4 more)

### Community 150 - "Config"
Cohesion: 0.06
Nodes (55): newCatalogRefresher(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted() (+47 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.09
Nodes (33): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), startSelfGuard(), displayEntries(), resolveGuardConfig(), runGuard(), runGuardHeadless() (+25 more)

### Community 155 - "safeio.go"
Cohesion: 0.25
Nodes (13): invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), removeCheckedFile(), AtomicWriteAt(), CreateExclNoFollow(), descend() (+5 more)

### Community 156 - "deploy"
Cohesion: 0.24
Nodes (15): applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), downloadAndVerify(), downloadReleaseFiles(), fetchGUI() (+7 more)

### Community 157 - "daemonPage"
Cohesion: 0.24
Nodes (3): daemonPage, daemonState, orDefault()

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.21
Nodes (20): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+12 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "dialog_test.go"
Cohesion: 0.36
Nodes (12): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+4 more)

### Community 164 - ".netIntegrityContainer"
Cohesion: 0.28
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (47): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+39 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (46): homePage, pageDef, pageID, reel, shell, snapEngines, snapNet, newLogo() (+38 more)

### Community 167 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 168 - "time.Time"
Cohesion: 0.14
Nodes (9): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID (+1 more)

### Community 169 - "io.Writer"
Cohesion: 0.16
Nodes (13): formatElapsed(), orNone(), printBinary(), printRecording(), candidates(), confirmable(), confirmAll(), describe() (+5 more)

### Community 170 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 171 - "runLockdown"
Cohesion: 0.12
Nodes (18): CheckBPFLSM(), configureDaemonLogging(), confirmMasterKeyOverwrite(), resolveConfigPath(), runBPFCheck(), runGenKey(), runLockdown(), runOneShotMode() (+10 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "runBackend"
Cohesion: 0.36
Nodes (6): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive()

### Community 175 - "runMonitor"
Cohesion: 0.11
Nodes (19): CheckEBPF(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), newPinGeneration(), prepareDaemonStart(), validateForwardEvents() (+11 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 178 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 179 - "main_test.go"
Cohesion: 0.20
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.23
Nodes (15): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+7 more)

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

### Community 194 - "finishBinary"
Cohesion: 0.27
Nodes (8): beyondEvents(), finishBinary(), isShellHelper(), noteWhitelisted(), TestFileReadersStartDeselected(), TestNoteWhitelisted(), TestShellHelpersStartDeselected(), exeVerdict

### Community 195 - "run.go"
Cohesion: 0.32
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "plan_test.go"
Cohesion: 0.09
Nodes (30): TestGeneratedLibraryBlockRoundTrips(), CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), TestRuntimeClassMarkers() (+22 more)

### Community 199 - "FetchLatest"
Cohesion: 0.24
Nodes (11): installHooks(), staged(), FetchLatest(), fetchRelease(), fetchReleases(), FetchVerified(), Release, TestFetchReleaseByTag() (+3 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 202 - "DaemonEvent"
Cohesion: 0.12
Nodes (27): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), newDaemonModel(), runDaemon(), runDaemonUI(), runHeadless() (+19 more)

### Community 203 - "Use"
Cohesion: 0.22
Nodes (8): channelOf(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestRunWithRefusesConcurrentRun(), cancelUI, TestEditTextMapsUICancel(), Use()

### Community 204 - "mountTable"
Cohesion: 0.39
Nodes (5): MountInfo, MountByID(), ParseMountinfo(), ReadMountinfo(), mountTable

### Community 205 - "uninstall/uninstall.go"
Cohesion: 0.09
Nodes (39): deletePostBackups(), listBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata(), confirmKeyDeletion(), decryptDirectories(), decryptStep() (+31 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - "guardModel"
Cohesion: 0.06
Nodes (34): listenForDaemonEvents(), tickDaemonStats(), formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), appendCapped(), clipToWidth() (+26 more)

### Community 210 - "loadDaemonConfig"
Cohesion: 0.25
Nodes (8): loadDaemonConfig(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds(), SetUserHomes(), TestOpenSystemPlacedRefusesUserHome()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 213 - "inode_dev"
Cohesion: 0.15
Nodes (18): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), current_is_inspector(), exe_tag_bits() (+10 more)

### Community 214 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 217 - "watchPattern"
Cohesion: 0.17
Nodes (15): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), systemBinaryPatterns(), systemPatterns() (+7 more)

### Community 218 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 219 - "run"
Cohesion: 0.36
Nodes (8): InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestReadInstalledVersion(), GUIInstalled()

### Community 220 - "shQuote"
Cohesion: 0.17
Nodes (5): IntegrationSuite, IntegrationSuite, overlayPlanter(), shQuote(), IntegrationSuite

### Community 221 - "NetEventType"
Cohesion: 0.08
Nodes (17): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+9 more)

### Community 222 - "Active"
Cohesion: 0.28
Nodes (8): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), EditText(), Active(), Stage()

### Community 223 - "selfGuards"
Cohesion: 0.52
Nodes (3): newSelfGuards(), startSelfGuards(), selfGuards

### Community 224 - "install/install.go"
Cohesion: 0.06
Nodes (65): anyDaemonRunning(), offerDaemonStart(), applyLiveRefresh(), asRoot(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled() (+57 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "securewrite_test.go"
Cohesion: 0.43
Nodes (6): refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeWithin()

### Community 227 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 228 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 229 - "links_test.go"
Cohesion: 0.54
Nodes (7): newLinkView(), checkLink(), linkFixture(), TestCheckLinksAccepts(), TestCheckLinksJudgesEachTarget(), TestCheckLinksRefusesUserFiles(), TestCheckLinksSeesThePlannedTree()

### Community 231 - "Prereq"
Cohesion: 0.33
Nodes (3): confirmRunPrereq(), Prereq, fakeVault

### Community 232 - "Existing"
Cohesion: 0.29
Nodes (3): sectionFor(), Existing, ExistingRule

### Community 234 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "linkView"
Cohesion: 0.43
Nodes (3): linkTarget(), entryKind, linkView

### Community 239 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 241 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 244 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 633 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **44 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.netIntegrityContainer`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.12241477888842886 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `sync.Mutex`, `GuardInodeKey`, `time.Duration`, `controlServer`, `BinaryEntry`, `guardUnitTest`, `time.Time`, `os.File`, `GuardEvent`, `daemonconfig.go`, `Vault`, `Resource`, `buildOneGuard`, `selfGuards`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02075055187637969 - nodes in this community are weakly interconnected._
- **Why does `backend` connect `Message` to `time.Duration`, `liveSession`, `runBackend`, `go_pkg_path_filepath`, `recordFlow`, `TestCleanRel`, `vaultfs.go`, `Client`, `Resource`, `Vault`, `sync.Mutex`, `backend_test.go`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._