# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 551 files · ~1,774,906 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7152 nodes · 32165 edges · 271 communities (211 shown, 60 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2928 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9beea37b`
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
- Candidate
- expandPlaceholders
- SanitizeText
- Guard
- Network
- startCatalogRefresh
- IntegrationSuite
- guard_trust.bpf.c
- fadescroll.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- Finder
- guard.bpf.c
- GuardInodeKey
- changelog.go
- Config
- LaunchRule
- OpenConfined
- NetGuardEvent
- policyEditor
- BinaryEntry
- picker
- stampCopyMetadata
- IntegrationSuite
- Message
- controlServer
- render.go
- pinstate.go
- IntegrationSuite
- absPath
- monitorUnitTest
- main_test.go
- fileedit_model_test.go
- fakeGuardRepo
- stdio
- update_test.go
- fadeScroll
- daemon_multicall_test.go
- editView
- Client
- watchGroup
- watchSet
- .apply
- githubRelease
- mutedText
- engine
- Monitor
- installPage
- Recorder
- exe_supersede.h
- filevault_test.go
- .showSized
- unix_connect.c
- wm.go
- IntegrationSuite
- github.com/charmbracelet/bubbletea.Cmd
- liveSession
- catalogRefresher
- RefreshSection
- golang.org/x/sys/unix.Stat_t
- fcntl
- bpfstats/main.go
- binaryVetter
- codecheck/main.go
- GlobReservations
- EventType
- process_vm_readv.c
- os.File
- op
- planUpdaters
- rcUser
- vaultfs.go
- journal
- dialLiveSession
- runDaemon
- GenerateConf
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- recoverFileInPlace
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- find
- runMonitor
- Well-known bypass classes
- buildTrustedSet
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- runNetworkGuard
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- deploy
- ConfinedEntry
- RunWith
- runBackend
- IntegrationSuite
- runWizard
- github.com/spf13/cobra.Command
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- parseDaemonEvents
- maintainerKeeper
- WebSocket /ws client
- inet_connect.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- uninstallAll
- .FilterExistingWhitelist
- trace-app-libs.sh
- IntegrationSuite
- pickfs.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- find_test.go
- serve.go
- NetEventType
- dialog_test.go
- Ledger
- maintainer/plan.go
- github.com/stretchr/testify/suite.Suite
- RootPlacedName
- runNetworkMonitor
- fakeStore
- loginServer
- IntegrationSuite
- buildOneGuard
- safeio.go
- copyRegularAt
- cursor.go
- Installation
- backend_test.go
- graphify-refresh.sh
- .netIntegrityContainer
- time.Time
- newTestTarget
- lockAndDeprovision
- sync.Mutex
- shell
- layout.go
- logGrid
- KernelDev
- SelectOrphans
- Limitations
- net_tester/main.go
- IntegrationSuite
- TrustGuard
- IntegrationSuite
- run.go
- Daemon configuration
- github.com/charmbracelet/bubbles/textarea.Model
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
- recordFlow
- .populateWatchInodes
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- SampleContent
- pathCache
- daemon/ports.go
- DaemonEvent
- io.Writer
- FileTarget
- secureResources
- Desktop GUI
- IntegrationSuite
- confirmUpdate
- GuardEvent
- .eventTypeAllowed
- ValidateTargets
- tempGrantSetup
- shQuote
- .awaitDaemonUpOrDead
- loadDaemonConfig
- auditGroupCoverage
- .TestRecord_GuardParity
- configedit_test.go
- selfGuards
- run
- prepareDaemonStart
- IntegrationSuite
- fakeDaemon
- supersede_probe.c
- procstats.go
- decryptStep
- ParseMountinfo
- bleedLayout
- github.com/Virgula0/app-listener/tools/codecheck
- refuseOutsideHardlinks
- Vault
- inode_dev.h
- PlanDaemonUnitUpdate
- writeWithin
- Troubleshooting
- ReplaceInstalledBinary
- SecretStoreNotes
- sources
- NewEventFanout
- guardLSMHooks
- checkRunningBinaryMatchesInstalled
- ReadInstallCert
- monitor.md
- launchKey
- User
- TestCleanRel
- bpfLsmListed
- prctl_gate
- mountTable
- parseKernelRelease
- golang.org/x/sys/unix.Timespec
- fileID
- MulticallError

## God Nodes (most connected - your core abstractions)
1. `Config` - 137 edges
2. `Guard` - 125 edges
3. `shell` - 97 edges
4. `Load()` - 95 edges
5. `absPath()` - 93 edges
6. `IntegrationSuite` - 83 edges
7. `fileEditModel` - 77 edges
8. `writeConfig()` - 74 edges
9. `GuardInodeKey` - 70 edges
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

## Communities (271 total, 60 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.09
Nodes (18): tempJournalRow, Progress, User, invalidPassword, recordTick, linkedFile, launchRule, trustEvent (+10 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (173): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod() (+165 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (7): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "Vault"
Cohesion: 0.14
Nodes (17): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+9 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (81): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+73 more)

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
Nodes (22): addErrKind, ReplacementCheck, classifyAddErr(), addRecorder(), copySelf(), guardUnitTest, runTool(), walkEntries() (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (34): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+26 more)

### Community 10 - "Candidate"
Cohesion: 0.06
Nodes (42): appendSectionsAndEdit(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+34 more)

### Community 11 - "expandPlaceholders"
Cohesion: 0.10
Nodes (14): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), symlinkTargetGlob(), expandPlaceholders(), BinaryRule, TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing() (+6 more)

### Community 12 - "SanitizeText"
Cohesion: 0.05
Nodes (25): backendLink, checkGrid, configView, forwardForm, forwardView, protectedPage, recordReview, recordView (+17 more)

### Community 13 - "Guard"
Cohesion: 0.08
Nodes (8): exeRow, rootHandle, supersededBinary, openRoot(), SharedPinDegraded(), TestNoteRootMountFollowsMountID(), Guard, releaseRootTrust()

### Community 14 - "Network"
Cohesion: 0.09
Nodes (10): binaryList, newBinaryList(), netGuardPill(), NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy (+2 more)

### Community 15 - "startCatalogRefresh"
Cohesion: 0.09
Nodes (25): backingDeviceUnion(), buildGuards(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded(), startCatalogRefresh(), startGuardedDaemon(), startGuardedDaemonAbortable() (+17 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (61): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+53 more)

### Community 18 - "fadescroll.go"
Cohesion: 0.03
Nodes (29): appTheme, binaryRow, clickPanel, homePage, iconTile, iconTileRenderer, navItem, navItemRenderer (+21 more)

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
Cohesion: 0.07
Nodes (10): Resource, GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, runtimeReleased(), uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (39): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+31 more)

### Community 24 - "Finder"
Cohesion: 0.16
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (90): inode_dev(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+82 more)

### Community 26 - "GuardInodeKey"
Cohesion: 0.06
Nodes (28): refusedReplacements, supersedeUnnamed(), appletTag(), buildOf(), ExeKey(), GuardInodeKey, IsMulticallRefusal(), mcNameKey() (+20 more)

### Community 27 - "changelog.go"
Cohesion: 0.12
Nodes (29): newChangelogModel(), showChangelog(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation() (+21 more)

### Community 28 - "Config"
Cohesion: 0.09
Nodes (36): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), launchRulesByPath(), TestLaunchRulesByPath(), startMaintainers(), startTrustGuard(), trustStartupError(), askBunTmpdirUsers() (+28 more)

### Community 29 - "LaunchRule"
Cohesion: 0.13
Nodes (18): allowedEnv(), binaryLaunchRules(), TestAllowedEnvKeepsOnlyRootPlaced(), argvEntry, applyLaunchRule(), envName(), LaunchRule, ParseAllowEnv() (+10 more)

### Community 30 - "OpenConfined"
Cohesion: 0.10
Nodes (33): lines(), pickBinaries(), pickLibs(), reviewRecording(), loaderName(), LoaderRefusal(), soNames(), MulticallRefusal() (+25 more)

### Community 31 - "NetGuardEvent"
Cohesion: 0.13
Nodes (14): snapNetGuard, NetGuardEvent, NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError() (+6 more)

### Community 32 - "policyEditor"
Cohesion: 0.13
Nodes (13): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+5 more)

### Community 33 - "BinaryEntry"
Cohesion: 0.06
Nodes (21): deferredBinary, keptEntries, BinaryEntry, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), checkBtrfsKeys(), canonicalBinaryPath() (+13 more)

### Community 34 - "picker"
Cohesion: 0.14
Nodes (6): keyEntry, picker, errText(), isFindShortcut(), newKeyEntry(), newPickRow()

### Community 35 - "stampCopyMetadata"
Cohesion: 0.40
Nodes (5): stampCopyMetadata(), TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 37 - "Message"
Cohesion: 0.09
Nodes (14): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+6 more)

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (23): readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+15 more)

### Community 39 - "render.go"
Cohesion: 0.12
Nodes (13): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+5 more)

### Community 40 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (5): IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "main_test.go"
Cohesion: 0.18
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "fakeGuardRepo"
Cohesion: 0.07
Nodes (15): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant(), newExeTap() (+7 more)

### Community 48 - "update_test.go"
Cohesion: 0.08
Nodes (25): downloadFile(), downloadReleaseFiles(), fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestDownloadFileHTTPError() (+17 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (28): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+20 more)

### Community 50 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 51 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 52 - "Client"
Cohesion: 0.10
Nodes (12): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+4 more)

### Community 53 - "watchGroup"
Cohesion: 0.08
Nodes (47): watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector() (+39 more)

### Community 54 - "watchSet"
Cohesion: 0.18
Nodes (11): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern (+3 more)

### Community 55 - ".apply"
Cohesion: 0.10
Nodes (12): Prereq, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier, fakeVault (+4 more)

### Community 56 - "githubRelease"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 57 - "mutedText"
Cohesion: 0.08
Nodes (41): filesPage, multiAnswer, recordForm, newLogo(), newBinaryRow(), newCheck(), newRadioGroup(), newSelect() (+33 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (40): engine, inodeBatchFn, portOwner, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel() (+32 more)

### Community 59 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 60 - "installPage"
Cohesion: 0.07
Nodes (16): actionGate, binSource, handCheck, installAction, installedBinary, installPage, latestRelease, uninstallCard (+8 more)

### Community 61 - "Recorder"
Cohesion: 0.16
Nodes (15): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+10 more)

### Community 63 - "filevault_test.go"
Cohesion: 0.19
Nodes (25): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsEncryptedDispatchesToFileVaultForRegularFiles() (+17 more)

### Community 64 - ".showSized"
Cohesion: 0.07
Nodes (23): InstallMode, InstallRequest, installUI, UninstallHooks, UninstallRequest, UpdateHooks, ask(), promptBody() (+15 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 68 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (52): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), formatGuardEventLine(), sanitizeTerminalText() (+44 more)

### Community 69 - "liveSession"
Cohesion: 0.16
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+5 more)

### Community 70 - "catalogRefresher"
Cohesion: 0.15
Nodes (9): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), catalogRefresher (+1 more)

### Community 71 - "RefreshSection"
Cohesion: 0.09
Nodes (29): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+21 more)

### Community 72 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.12
Nodes (22): placedWalk, treeJudge, mountinfoDev(), OnBtrfs(), SuperblockDev(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs(), confineBelow() (+14 more)

### Community 73 - "fcntl"
Cohesion: 0.08
Nodes (11): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), hold_forever() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.14
Nodes (23): CheckBPFLSM(), runBPFCheck(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity() (+15 more)

### Community 75 - "binaryVetter"
Cohesion: 0.07
Nodes (21): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), systemPatterns(), systemRule(), newHashCache() (+13 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "GlobReservations"
Cohesion: 0.12
Nodes (16): reserveChain(), GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, ParseGlobName(), resolveBits() (+8 more)

### Community 78 - "EventType"
Cohesion: 0.09
Nodes (12): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, Cstr(), DecodeBpfEvent(), MonitorRepository (+4 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.08
Nodes (37): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp (+29 more)

### Community 81 - "op"
Cohesion: 0.17
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 82 - "planUpdaters"
Cohesion: 0.04
Nodes (51): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+43 more)

### Community 83 - "rcUser"
Cohesion: 0.09
Nodes (54): setupBunTmpdir(), installFacultatives(), installSSHAgentFacultative(), listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf() (+46 more)

### Community 84 - "vaultfs.go"
Cohesion: 0.09
Nodes (25): DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode(), inodeOf() (+17 more)

### Community 85 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 86 - "dialLiveSession"
Cohesion: 0.07
Nodes (33): auditAfterEdit(), auditFindings(), TestValidatePassword(), ValidatePassword(), dialLiveSession(), applyEdit(), checkAndApply(), editAgain() (+25 more)

### Community 87 - "runDaemon"
Cohesion: 0.13
Nodes (11): newControlManager(), startControlManager(), catchLifecycleSignals(), markReady(), markReadyWithLogin(), notifySystemdReady(), runDaemon(), startPprof() (+3 more)

### Community 88 - "GenerateConf"
Cohesion: 0.05
Nodes (62): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+54 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.08
Nodes (10): daemonPage, daemonState, runHost, runView, stageItem, stageState, journalLine(), TestJournalLine() (+2 more)

### Community 92 - "recoverFileInPlace"
Cohesion: 0.19
Nodes (17): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+9 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.07
Nodes (22): foreignPinnedLink(), FirstOnBtrfs(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), codeVouched(), isTraced() (+14 more)

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

### Community 100 - "find"
Cohesion: 0.22
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 101 - "runMonitor"
Cohesion: 0.19
Nodes (14): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+6 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "runNetworkGuard"
Cohesion: 0.12
Nodes (20): DrainUntilSignal(), runGuardHeadless(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard() (+12 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 110 - "deploy"
Cohesion: 0.14
Nodes (26): installHooks(), staged(), archAssets(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), deploy(), downloadAndVerify() (+18 more)

### Community 111 - "ConfinedEntry"
Cohesion: 0.11
Nodes (32): binaryVerifyState, HashStore, sharedHashEntry, VettedInode, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared() (+24 more)

### Community 112 - "RunWith"
Cohesion: 0.14
Nodes (17): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), stagedBinary, Mode, Request (+9 more)

### Community 113 - "runBackend"
Cohesion: 0.27
Nodes (8): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), anyDaemonRunning(), changePasswordLive()

### Community 115 - "runWizard"
Cohesion: 0.06
Nodes (46): offerDaemonStart(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), deploy(), installBinaryOnly(), markInstallComplete(), runRefreshCatalogOnly() (+38 more)

### Community 116 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (39): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe() (+31 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "parseDaemonEvents"
Cohesion: 0.20
Nodes (5): IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 121 - "maintainerKeeper"
Cohesion: 0.20
Nodes (10): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp, maintainerKeeper (+2 more)

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

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "uninstallAll"
Cohesion: 0.12
Nodes (19): releaseBuild(), Run(), confirmKeyDeletion(), run(), runUninstall(), RunWith(), removeKeyAndEmptyDir(), removeMasterKey() (+11 more)

### Community 131 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "pickfs.go"
Cohesion: 0.29
Nodes (12): listOpts, pickEntry, treeSearch, plural(), entryIsDir(), filterEntries(), listDir(), matchesTerms() (+4 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 143 - "serve.go"
Cohesion: 0.08
Nodes (31): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+23 more)

### Community 144 - "NetEventType"
Cohesion: 0.07
Nodes (21): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+13 more)

### Community 145 - "dialog_test.go"
Cohesion: 0.36
Nodes (12): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+4 more)

### Community 146 - "Ledger"
Cohesion: 0.05
Nodes (46): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+38 more)

### Community 147 - "maintainer/plan.go"
Cohesion: 0.18
Nodes (21): isWhiteout(), opaque(), pathByte(), pathRecord(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent() (+13 more)

### Community 148 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.19
Nodes (5): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option, editorHarness

### Community 149 - "RootPlacedName"
Cohesion: 0.16
Nodes (14): RootPlacedName(), rootOnlyTree(), TestRootPlacedName(), userAssignable(), linkTarget(), newLinkView(), checkLink(), linkFixture() (+6 more)

### Community 150 - "runNetworkMonitor"
Cohesion: 0.11
Nodes (11): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runNetworkMonitor(), runTUI(), NetworkMonitorRepository, NewNetModel() (+3 more)

### Community 152 - "loginServer"
Cohesion: 0.19
Nodes (10): hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor(), startLoginServer(), loginServer, madeDir (+2 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.08
Nodes (37): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), startSelfGuard(), displayEntries(), resolveGuardConfig(), runGuard(), runGuardTUI() (+29 more)

### Community 155 - "safeio.go"
Cohesion: 0.10
Nodes (28): readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), ancestorEnv(), invokingUID() (+20 more)

### Community 156 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - "cursor.go"
Cohesion: 0.08
Nodes (19): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+11 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - ".netIntegrityContainer"
Cohesion: 0.28
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 162 - "time.Time"
Cohesion: 0.21
Nodes (7): LockedError, newWatchSet(), fakeInotify, sweepID, sweepMemo, Progress, Seen

### Community 163 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 164 - "lockAndDeprovision"
Cohesion: 0.12
Nodes (11): entryOf(), deprovisionKind, migrationTarget, writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr(), classifyMigrationTarget(), isDeprovisionBusy() (+3 more)

### Community 165 - "sync.Mutex"
Cohesion: 0.05
Nodes (48): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+40 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (58): fakeLink, Options, pageDef, pageID, reel, resizeWatch, shell, sizedDialog (+50 more)

### Community 167 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 168 - "logGrid"
Cohesion: 0.10
Nodes (10): logGrid, logHook, logLine, logRow(), logStyle(), newLogGrid(), formatEntry(), sortedKeys() (+2 more)

### Community 169 - "KernelDev"
Cohesion: 0.14
Nodes (20): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), binaryStatOf(), KernelDev(), StatBinary() (+12 more)

### Community 170 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 175 - "TrustGuard"
Cohesion: 0.05
Nodes (31): tempGrant, tempMaskOp, tempRow, trustRows, deleteInoKeys(), deleteResKeys(), droppableFreed(), freedSuperseded() (+23 more)

### Community 176 - "IntegrationSuite"
Cohesion: 0.24
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 177 - "run.go"
Cohesion: 0.36
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 178 - "Daemon configuration"
Cohesion: 0.18
Nodes (11): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Port sections, Reloading (+3 more)

### Community 179 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (22): class, pos, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+14 more)

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

### Community 191 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "recordFlow"
Cohesion: 0.13
Nodes (17): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+9 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "SampleContent"
Cohesion: 0.09
Nodes (28): deployFiles(), grantNewBinaries(), preflightDeployedBinary(), TestDaemonUnitWithMetadataOutput(), dupFD(), installConfig(), installFileAs(), installServices() (+20 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "daemon/ports.go"
Cohesion: 0.36
Nodes (7): ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), parseReservedPorts(), readReservedPorts(), TestParseReservedPorts()

### Community 202 - "DaemonEvent"
Cohesion: 0.12
Nodes (30): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+22 more)

### Community 203 - "io.Writer"
Cohesion: 0.12
Nodes (11): formatElapsed(), orNone(), printBinary(), printRecording(), TestProgressReader(), fyneLogFilter, bottomBar, newBottomBar() (+3 more)

### Community 204 - "FileTarget"
Cohesion: 0.10
Nodes (16): specFromForm(), guardPill(), newFilesPage(), specSummary(), EventTypes(), maskFor(), Files, FileSpec (+8 more)

### Community 205 - "secureResources"
Cohesion: 0.24
Nodes (8): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 206 - "Desktop GUI"
Cohesion: 0.13
Nodes (15): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+7 more)

### Community 209 - "confirmUpdate"
Cohesion: 0.29
Nodes (7): confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), Use()

### Community 210 - "GuardEvent"
Cohesion: 0.08
Nodes (19): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), TestLocalizePaths() (+11 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "ValidateTargets"
Cohesion: 0.14
Nodes (12): dedupe(), existingSection(), Plan, guardedMatch(), libOfChosen(), overlaps(), refusedPath(), sectionFor() (+4 more)

### Community 213 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 214 - "shQuote"
Cohesion: 0.17
Nodes (5): IntegrationSuite, IntegrationSuite, overlayPlanter(), shQuote(), IntegrationSuite

### Community 216 - "loadDaemonConfig"
Cohesion: 0.20
Nodes (10): loadDaemonConfig(), noteReloadedSettings(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds(), SetUserHomes() (+2 more)

### Community 217 - "auditGroupCoverage"
Cohesion: 0.22
Nodes (11): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+3 more)

### Community 218 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 219 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 220 - "selfGuards"
Cohesion: 0.33
Nodes (5): newSelfGuards(), selfProtectSpecs(), startSelfGuards(), selfGuards, selfProtectSpec

### Community 221 - "run"
Cohesion: 0.18
Nodes (14): channelOf(), fetchRelease(), InstalledGUIVersion(), InstalledVersion(), Latest(), readInstalledVersion(), run(), runUpdate() (+6 more)

### Community 222 - "prepareDaemonStart"
Cohesion: 0.10
Nodes (21): CheckEBPF(), configureDaemonLogging(), confirmMasterKeyOverwrite(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runGenKey(), runLockdown() (+13 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 227 - "decryptStep"
Cohesion: 0.22
Nodes (9): askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), TestExistingEncryptedDirsNilPaths(), TestExistingEncryptedDirsSkipsPlain(), ExistingEncryptedDirs() (+1 more)

### Community 228 - "ParseMountinfo"
Cohesion: 0.12
Nodes (15): MountInfo, mountinfoUnitTest, hostVouchedDevs(), isFuseType(), MountByID(), parseMountDev(), ParseMountinfo(), PathUnder() (+7 more)

### Community 231 - "refuseOutsideHardlinks"
Cohesion: 0.28
Nodes (8): lstatOf(), outsideHardlinks(), refuseOutsideHardlinks(), TestRefuseOutsideHardlinks_FileRoot(), TestRefuseOutsideHardlinks_InTreeLinksAccepted(), TestRefuseOutsideHardlinks_MissingRootFails(), TestRefuseOutsideHardlinks_OutsideLinkRefused(), writeFile()

### Community 233 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 234 - "PlanDaemonUnitUpdate"
Cohesion: 0.31
Nodes (9): ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate(), TestPlanDaemonUnitUpdateDropsLocalEdits(), TestShippedDaemonUnit() (+1 more)

### Community 235 - "writeWithin"
Cohesion: 0.25
Nodes (8): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "ReplaceInstalledBinary"
Cohesion: 0.25
Nodes (8): installGUI(), ensureInstalledBinary(), prepareInstallation(), TestEnsureInstalledBinary(), deployGUI(), DeployGUIBinary(), ReplaceInstalledBinary(), TestReplaceInstalledBinary()

### Community 238 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 240 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 241 - "guardLSMHooks"
Cohesion: 0.25
Nodes (5): lsmHook, trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), TestSocketHooksAreRequired()

### Community 242 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 243 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 244 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 246 - "User"
Cohesion: 0.10
Nodes (29): entryWriters(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), bunLaunchersOffered(), removeFacultative() (+21 more)

### Community 248 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 249 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 251 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

## Knowledge Gaps
- **33 isolated node(s):** `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 678 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **60 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `.netIntegrityContainer`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `IntegrationSuite`, `.TestRecord_GuardParity`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.0874450604916989 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `BinaryEntry`, `fakeDaemon`, `time.Duration`, `engine`, `sync.Mutex`, `time.Time`, `guardUnitTest`, `DaemonEvent`, `binaryVetter`, `startCatalogRefresh`, `TrustGuard`, `os.File`, `ConfinedEntry`, `GuardEvent`, `runDaemon`, `buildOneGuard`, `selfGuards`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.016571632216678546 - nodes in this community are weakly interconnected._
- **Why does `shell` connect `shell` to `.showSized`, `policyEditor`, `picker`, `time.Duration`, `wm.go`, `sync.Mutex`, `go_pkg_strings`, `SanitizeText`, `FileTarget`, `Network`, `dialog_test.go`, `fadescroll.go`, `mutedText`, `runView`, `installPage`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09441441441441441 - nodes in this community are weakly interconnected._