# Graph Report - app-listener  (2026-10-19)

## Corpus Check
- 499 files · ~1,723,172 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6582 nodes · 28692 edges · 254 communities (206 shown, 48 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2668 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e8e6d5ce`
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
- mkdirs
- GuardEvent
- FileTarget
- protectedPage
- fileedit_symlink_test.go
- runNetworkGuard
- os.FileMode
- IntegrationSuite
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- pinstate.go
- guard.bpf.c
- TrustGuard
- changelog.go
- highlight_test.go
- rcUser
- github.com/charmbracelet/bubbles/textarea.Model
- TestMonitorUseCaseLifecycle
- UI
- github.com/charmbracelet/bubbletea.Cmd
- picker
- runOfflineEdit
- IntegrationSuite
- Message
- controlServer
- EventType
- NetEventType
- IntegrationSuite
- absPath
- monitorUnitTest
- RunWith
- fileedit_model_test.go
- CandidateDir
- CopyTreeWithProgress
- fadeScroll
- binaryList
- cursor.go
- Client
- watchGroup
- fyne.io/fyne/v2/widget.Label
- Plan
- FileEvent
- mutedText
- engine
- os.File
- filevault_test.go
- Recorder
- exe_supersede.h
- installPage
- catalogRefresher
- sync.Mutex
- gui/session.go
- IntegrationSuite
- dialLiveSession
- update_test.go
- Network
- .FilterExistingWhitelist
- Guard
- fcntl
- bpfstats/main.go
- daemon_multicall_test.go
- RefreshSection
- journal
- time.Time
- process_vm_readv.c
- compileGlobNames
- editView
- planUpdaters
- update/update.go
- net_tester/main.go
- PlanDaemonUnitUpdate
- recordFlow
- go_pkg_github_com_testcontainers_testcontainers_go
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- ResolvePinBase
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- ConfinedEntry
- render.go
- Well-known bypass classes
- systemd.go
- jit_provenance.c
- inspector_probe.c
- recordops.c
- IntegrationSuite
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- find
- mountTable
- guardModel
- .step
- IntegrationSuite
- io.Writer
- .netIntegrityContainer
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
- newTestTarget
- updatePage
- trace-app-libs.sh
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- NetGuardEvent
- serveWebSocket
- run
- Open
- Ledger
- buildTrustedSet
- filevault.go
- SelectOrphans
- runMonitor
- fakeStore
- inodeChain
- IntegrationSuite
- SanitizeText
- NetworkMonitor
- NewDumpHook
- hashCache
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- vaultfs.go
- .populateWatchInodes
- feed.go
- shell
- runGenKey
- auditGroupCoverage
- runDaemon
- WithBottomBar
- logLine
- Limitations
- Troubleshooting
- IntegrationSuite
- TestMain
- netinject.c
- copyRegularAt
- inode_dev.h
- Binary
- Compatibility
- How it works
- github.com/spf13/cobra.Command
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- changelog_test.go
- writeWithin
- edit-protected
- monitor.md
- kernel-prof.sh
- deploy
- pathCache
- logGrid
- DaemonEvent
- AddServeFlags
- ValidateTargets
- UnifiedDiff
- Desktop GUI
- IntegrationSuite
- Config
- run
- .eventTypeAllowed
- bleedLayout
- bpfLsmListed
- IntegrationSuite
- .awaitDaemonUpOrDead
- LibraryClosure
- KernelDev
- launchKey
- NewEventFanout
- TestCleanRel
- formatEntry
- prctl_gate
- ReadInstallCert
- tempGrantSetup
- supersede_probe.c
- Prereq
- ptrace_race.c
- downloadFile
- Vault
- resolveFilesystemPrereqs
- copyXattrs
- Execute
- taint_victim.c
- .runEditorHarness
- parseKernelRelease
- golang.org/x/sys/unix.Timespec
- serve.go
- configedit_test.go
- MulticallError
- RawTarget
- MulticallNameTags
- IntegrationSuite
- Daemon configuration
- serveHandler
- parseGuardEvent
- WithSuspectsFrom
- loadDaemonConfig
- github.com/stretchr/testify/suite.Suite
- checkRunningBinaryMatchesInstalled
- unlockUnderGuard
- parseResize
- bytes.Buffer

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 109 edges
3. `GuardInodeKey` - 99 edges
4. `shell` - 91 edges
5. `Load()` - 89 edges
6. `absPath()` - 88 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

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

## Communities (254 total, 48 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (18): tempJournalRow, Progress, User, invalidPassword, recordTick, linkedFile, launchRule, trustEvent (+10 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (151): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestEveryInstallModeMapped() (+143 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.12
Nodes (21): verifyEncryptionState(), deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy() (+13 more)

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
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (25): addErrKind, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+17 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "GuardEvent"
Cohesion: 0.04
Nodes (16): holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+8 more)

### Community 11 - "FileTarget"
Cohesion: 0.07
Nodes (22): netGuardTab, policyEditor, policyHooks, policyTab, guardPill(), newPolicyTab(), verbFor(), netWarnings() (+14 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (21): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, recordForm, recordReview (+13 more)

### Community 13 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 14 - "runNetworkGuard"
Cohesion: 0.15
Nodes (15): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), DiscoverInfraBinaries(), Mode (+7 more)

### Community 15 - "os.FileMode"
Cohesion: 0.11
Nodes (24): TestDaemonUnitWithMetadataOutput(), dupFD(), installFile(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow() (+16 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "runView"
Cohesion: 0.09
Nodes (8): daemonPage, daemonState, runHost, runView, stageItem, stageState, captureLogs(), orDefault()

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
Cohesion: 0.06
Nodes (12): ownerPaths(), fakeDaemon, findResource(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase (+4 more)

### Community 23 - "auth_test.go"
Cohesion: 0.12
Nodes (41): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+33 more)

### Community 24 - "pinstate.go"
Cohesion: 0.13
Nodes (28): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder() (+20 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+64 more)

### Community 26 - "TrustGuard"
Cohesion: 0.05
Nodes (30): setConfinedHomes(), tempGrant, tempMaskOp, tempRow, trustHook, trustRows, ensureBtrfsLayout(), appendInodeRow() (+22 more)

### Community 27 - "changelog.go"
Cohesion: 0.11
Nodes (27): findButton, findAll(), lowerRunes(), NewFinder(), newFindInput(), overlay(), runesEqual(), findEditor() (+19 more)

### Community 28 - "highlight_test.go"
Cohesion: 0.13
Nodes (26): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+18 more)

### Community 29 - "rcUser"
Cohesion: 0.13
Nodes (34): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+26 more)

### Community 30 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.11
Nodes (25): class, pos, span, classOf(), column(), dedent(), Edit(), leadingIndent() (+17 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.27
Nodes (10): NewGuardUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkGuardUseCaseStartError() (+2 more)

### Community 32 - "UI"
Cohesion: 0.08
Nodes (22): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, ShowDiff() (+14 more)

### Community 33 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (11): recordTickCmd(), recordModel, diffModel, editorModel, renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel (+3 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (31): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, lstatOf(), outsideHardlinks() (+23 more)

### Community 35 - "runOfflineEdit"
Cohesion: 0.10
Nodes (19): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt() (+11 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (10): checkSecret(), File, backend, backendOp, backendState, fakeLink, BinaryCheck, Message (+2 more)

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (22): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+14 more)

### Community 39 - "EventType"
Cohesion: 0.06
Nodes (34): EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), binaryStatOf(), TempBinary, TemporaryGrant, TempRule (+26 more)

### Community 40 - "NetEventType"
Cohesion: 0.07
Nodes (18): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+10 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (5): IntegrationSuite, IntegrationSuite, IntegrationSuite, IntegrationSuite, absPath()

### Community 44 - "RunWith"
Cohesion: 0.17
Nodes (14): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), Mode, Request (+6 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "CandidateDir"
Cohesion: 0.08
Nodes (23): reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), expandPlaceholders() (+15 more)

### Community 48 - "CopyTreeWithProgress"
Cohesion: 0.14
Nodes (10): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), copyTreeRoot(), CopyTreeWithProgress(), TestCopyTreeProgress() (+2 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (26): Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+18 more)

### Community 51 - "cursor.go"
Cohesion: 0.07
Nodes (23): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+15 more)

### Community 52 - "Client"
Cohesion: 0.08
Nodes (18): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+10 more)

### Community 53 - "watchGroup"
Cohesion: 0.10
Nodes (38): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+30 more)

### Community 54 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (24): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel (+16 more)

### Community 55 - "Plan"
Cohesion: 0.12
Nodes (15): BinaryRule, InsertSectionsBeforeLibraries(), Plan, materialized(), TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+7 more)

### Community 56 - "FileEvent"
Cohesion: 0.11
Nodes (9): snapMonitor, BpfEvent, eventUnitTest, FileEvent, Cstr(), DecodeBpfEvent(), newFakeMonitorRepo(), fakeSource (+1 more)

### Community 57 - "mutedText"
Cohesion: 0.08
Nodes (38): filesPage, multiAnswer, newLogo(), newBinaryList(), newBinaryRow(), selectable(), newCheck(), TestAdaptiveGridDropsColumns() (+30 more)

### Community 58 - "engine"
Cohesion: 0.05
Nodes (29): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+21 more)

### Community 59 - "os.File"
Cohesion: 0.05
Nodes (65): VettedInode, Multicall, MulticallKind, loaderName(), LoaderRefusal(), soNames(), hostLoader(), TestLoaderRefusal() (+57 more)

### Community 60 - "filevault_test.go"
Cohesion: 0.18
Nodes (26): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestEnsureRecoverySidecarPlaceholder() (+18 more)

### Community 61 - "Recorder"
Cohesion: 0.16
Nodes (15): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), sortedKeys() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+7 more)

### Community 63 - "installPage"
Cohesion: 0.10
Nodes (16): actionGate, binSource, installAction, InstallMode, installPage, InstallRequest, Options, stagedBinary (+8 more)

### Community 64 - "catalogRefresher"
Cohesion: 0.15
Nodes (8): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 65 - "sync.Mutex"
Cohesion: 0.11
Nodes (14): newActivity(), newExeResolver(), guardRow(), TestRingKeepsNewest(), NewNetwork(), Counts, Filter, Row (+6 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (40): invokingIDs(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.26
Nodes (3): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite

### Community 68 - "dialLiveSession"
Cohesion: 0.09
Nodes (27): auditAfterEdit(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply(), editAgain() (+19 more)

### Community 69 - "update_test.go"
Cohesion: 0.09
Nodes (18): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+10 more)

### Community 70 - "Network"
Cohesion: 0.13
Nodes (8): netGuardPill(), NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy, Network, NetGuardConfig

### Community 71 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 72 - "Guard"
Cohesion: 0.04
Nodes (38): binaryVerifyState, deferredBinary, exeRow, keptEntries, refusedReplacements, rootHandle, BinaryEntry, openRoot() (+30 more)

### Community 73 - "fcntl"
Cohesion: 0.09
Nodes (11): main(), report(), denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 76 - "RefreshSection"
Cohesion: 0.11
Nodes (25): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+17 more)

### Community 77 - "journal"
Cohesion: 0.15
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 78 - "time.Time"
Cohesion: 0.07
Nodes (29): LockedError, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), isSymlink() (+21 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.19
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 81 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), newEditView(), relOf()

### Community 82 - "planUpdaters"
Cohesion: 0.15
Nodes (17): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+9 more)

### Community 83 - "update/update.go"
Cohesion: 0.23
Nodes (18): newReleaseWindow(), releasesSince(), channelOf(), compareStableVersions(), filterChannel(), findUpdate(), hasAsset(), Latest() (+10 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "PlanDaemonUnitUpdate"
Cohesion: 0.27
Nodes (10): ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate(), TestPlanDaemonUnitUpdateDropsLocalEdits(), TestShippedDaemonUnit() (+2 more)

### Community 86 - "recordFlow"
Cohesion: 0.14
Nodes (17): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+9 more)

### Community 87 - "go_pkg_github_com_testcontainers_testcontainers_go"
Cohesion: 0.16
Nodes (4): attrOpCase, inNosuidNamespace(), guardEventTypes(), guardExploitTest

### Community 88 - "liveSession"
Cohesion: 0.13
Nodes (15): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries() (+7 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 92 - "Monitor"
Cohesion: 0.16
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (14): foreignPinnedLink(), NewMountVouch(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs() (+6 more)

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
Nodes (14): procStats, statsSample, formatClock(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "ConfinedEntry"
Cohesion: 0.12
Nodes (29): HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats(), publishHash() (+21 more)

### Community 101 - "render.go"
Cohesion: 0.13
Nodes (13): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+5 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "systemd.go"
Cohesion: 0.11
Nodes (39): offerDaemonStart(), buildBinaryIfNeeded(), installBinaryOnly(), mustCwd(), revertSystemFiles(), readDaemonState(), DaemonUnitInstalled(), DaemonUnitState() (+31 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (14): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+6 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "find"
Cohesion: 0.22
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "guardModel"
Cohesion: 0.06
Nodes (35): runTUI(), Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero() (+27 more)

### Community 113 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 115 - "io.Writer"
Cohesion: 0.19
Nodes (12): formatElapsed(), orNone(), printBinary(), printRecording(), candidates(), confirmable(), confirmAll(), describe() (+4 more)

### Community 116 - ".netIntegrityContainer"
Cohesion: 0.28
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (60): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+52 more)

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

### Community 130 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 131 - "updatePage"
Cohesion: 0.20
Nodes (9): installedBinary, latestRelease, updatePage, hashFile(), hashReader(), readInstalled(), newUpdatePage(), boolIndex() (+1 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "NetGuardEvent"
Cohesion: 0.16
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 143 - "serveWebSocket"
Cohesion: 0.19
Nodes (10): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop() (+2 more)

### Community 144 - "run"
Cohesion: 0.12
Nodes (18): anyDaemonRunning(), confirmKeyDeletion(), run(), runUninstall(), RunWith(), removeKeyAndEmptyDir(), removeMasterKey(), revertBunLaunchers() (+10 more)

### Community 145 - "Open"
Cohesion: 0.16
Nodes (20): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+12 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "buildTrustedSet"
Cohesion: 0.14
Nodes (19): fileID(), openBunTmp(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure() (+11 more)

### Community 148 - "filevault.go"
Cohesion: 0.25
Nodes (13): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+5 more)

### Community 149 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 150 - "runMonitor"
Cohesion: 0.13
Nodes (11): MakeDisplayPaths(), ParseEventsFlag(), validateForwardEvents(), runHeadless(), runMonitor(), runTUI(), MonitorRepository, NewModel() (+3 more)

### Community 152 - "inodeChain"
Cohesion: 0.23
Nodes (15): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+7 more)

### Community 154 - "SanitizeText"
Cohesion: 0.04
Nodes (74): CheckBPFLSM(), CheckEBPF(), ParseNetEventsFlag(), multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), backingDeviceUnion(), buildConcurrency() (+66 more)

### Community 155 - "NetworkMonitor"
Cohesion: 0.18
Nodes (5): NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 156 - "NewDumpHook"
Cohesion: 0.19
Nodes (15): TestSetupDumpLogRefusesSymlink(), applyVerbosity(), effectiveVerbose(), preRun(), setupDumpLog(), VerboseLevel, BridgeStdLog(), NewDumpHook() (+7 more)

### Community 157 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "layout.go"
Cohesion: 0.24
Nodes (8): adaptiveGrid, checkGrid, gapVBox, handCheck, uninstallCard, cellMin(), rowHeights(), visible()

### Community 163 - "vaultfs.go"
Cohesion: 0.08
Nodes (29): checkFreshDir(), openDirNoFollow(), renewDir(), entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent() (+21 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (40): TestForwardFlagSurface(), actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction (+32 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (49): homePage, pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines (+41 more)

### Community 167 - "runGenKey"
Cohesion: 0.20
Nodes (11): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen() (+3 more)

### Community 168 - "auditGroupCoverage"
Cohesion: 0.20
Nodes (12): auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+4 more)

### Community 169 - "runDaemon"
Cohesion: 0.10
Nodes (20): newControlManager(), startControlManager(), catchLifecycleSignals(), makeReloadHandler(), markReady(), notifySystemdReady(), runDaemon(), startCatalogRefresh() (+12 more)

### Community 170 - "WithBottomBar"
Cohesion: 0.31
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 171 - "logLine"
Cohesion: 0.17
Nodes (7): logHook, logLine, journalLine(), TestJournalLine(), logRow(), logStyle(), DumpHook

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "copyRegularAt"
Cohesion: 0.36
Nodes (9): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), fchownFromInfo(), futimesFromInfo(), readDirNames(), readlinkFD() (+1 more)

### Community 178 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 179 - "Binary"
Cohesion: 0.13
Nodes (20): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+12 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.19
Nodes (14): credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+6 more)

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

### Community 194 - "changelog_test.go"
Cohesion: 0.22
Nodes (12): changelogText(), cumulativeChangelog(), newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize() (+4 more)

### Community 195 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "deploy"
Cohesion: 0.12
Nodes (28): installHooks(), releaseBuild(), Run(), staged(), applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions() (+20 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 202 - "DaemonEvent"
Cohesion: 0.07
Nodes (41): UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+33 more)

### Community 203 - "AddServeFlags"
Cohesion: 0.43
Nodes (7): AddHeadlessFlag(), AddServeFlags(), init(), init(), init(), init(), init()

### Community 204 - "ValidateTargets"
Cohesion: 0.18
Nodes (9): existingSection(), guardedMatch(), overlaps(), refusedPath(), validateTarget(), ValidateTargets(), TestValidateTargets(), Existing (+1 more)

### Community 205 - "UnifiedDiff"
Cohesion: 0.20
Nodes (11): applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), refreshCatalogConfig(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr(), installConfig(), ConfirmOverwrite() (+3 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - "Config"
Cohesion: 0.05
Nodes (76): entryWriters(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), addedResources(), appendSectionsAndEdit(), applyDiffAdditions() (+68 more)

### Community 210 - "run"
Cohesion: 0.24
Nodes (11): deployGUI(), InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), RunWith(), syncGUI() (+3 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 216 - "LibraryClosure"
Cohesion: 0.07
Nodes (27): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+19 more)

### Community 217 - "KernelDev"
Cohesion: 0.09
Nodes (31): BinaryStat, isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs(), binaryStatOf() (+23 more)

### Community 219 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 221 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 223 - "ReadInstallCert"
Cohesion: 0.21
Nodes (12): certifyInstall(), markInstallComplete(), InstallCert, InstallCertPath(), ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip() (+4 more)

### Community 224 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 225 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 227 - "ptrace_race.c"
Cohesion: 0.52
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 228 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 230 - "resolveFilesystemPrereqs"
Cohesion: 0.28
Nodes (6): encryptDirectories(), confirmRunPrereq(), resolveFilesystemPrereqs(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 231 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 232 - "Execute"
Cohesion: 0.67
Nodes (3): Execute(), initLogger(), main()

### Community 235 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 237 - "serve.go"
Cohesion: 0.22
Nodes (10): isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), Serve(), serveHTTP(), waitProgram() (+2 more)

### Community 238 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 241 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 242 - "MulticallNameTags"
Cohesion: 0.22
Nodes (6): appletTag(), buildOf(), mcNameKey(), MulticallNameTags(), TestAppletTagStablePerName(), TestMcNameKey()

### Community 244 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 245 - "serveHandler"
Cohesion: 0.31
Nodes (8): basicAuth(), securityHeaders(), serveHandler(), TestBasicAuth(), TestSecurityHeadersAndExactRoot(), TestServeHandlerProtectsPageAndWebSocket(), TestServeWebSocketSnapshotsResizeAndSingleViewer(), validateHost()

### Community 246 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 247 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 248 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), noteReloadedSettings(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 250 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 251 - "unlockUnderGuard"
Cohesion: 0.50
Nodes (5): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), SectionScan

### Community 252 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 604 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **48 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `liveSession`, `find`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0944115673235312 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01846327501410949 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `.netIntegrityContainer`, `IntegrationSuite`, `IntegrationSuite`, `go_pkg_github_com_testcontainers_testcontainers_go`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._