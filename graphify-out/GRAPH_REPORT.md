# Graph Report - app-listener  (2026-10-19)

## Corpus Check
- 486 files · ~1,715,324 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6481 nodes · 27998 edges · 233 communities (193 shown, 40 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2602 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6d6789b0`
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
- BinaryEntry
- policyEditor
- protectedPage
- codeedit.go
- CandidateDir
- editView
- IntegrationSuite
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- bpfstats/main.go
- guard.bpf.c
- TrustGuard
- Finder
- github.com/charmbracelet/bubbletea.Cmd
- rcUser
- Config
- MonitorRepository
- UI
- NetEventType
- picker
- fileedit_model_test.go
- IntegrationSuite
- Message
- controlServer
- TempBinary
- ConfinedEntry
- IntegrationSuite
- absPath
- monitorUnitTest
- runMaintenanceMode
- runDaemon
- User
- stdio
- lockAndDeprovision
- fadeScroll
- newTestTarget
- cursor.go
- Client
- watchGroup
- tabs.go
- Plan
- RawTarget
- mutedText
- GuardInodeKey
- os.File
- Candidate
- Recorder
- exe_supersede.h
- Network
- RefreshSection
- pinstate.go
- wm.go
- IntegrationSuite
- Binary
- FetchLatest
- fcntl
- .FilterExistingWhitelist
- Guard
- FileTarget
- os.FileMode
- catalogRefresher
- journal
- deploy
- process_vm_readv.c
- GlobReservations
- serveWebSocket
- planUpdaters
- githubRelease
- net_tester/main.go
- runBackend
- github.com/charmbracelet/bubbles/textarea.Model
- recordFlow
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- inodeChain
- Monitor
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- highlight_test.go
- Prereq
- Well-known bypass classes
- installPage
- jit_provenance.c
- inspector_probe.c
- recordops.c
- openFileVaultWithMasterKey
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- recoverFileInPlace
- mountTable
- NetGuardEvent
- .step
- IntegrationSuite
- readGrantRequest
- Vault
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
- .launch
- supersede_probe.c
- deploy
- trace-app-libs.sh
- IntegrationSuite
- dialLiveSession
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- GuardEvent
- batchInodeRows
- runOfflineEdit
- Open
- ValidateTargets
- install_run.go
- testKeyPair
- SelectOrphans
- TrustGlob
- fakeStore
- filevault_test.go
- IntegrationSuite
- buildOneGuard
- secureResources
- IntegrationSuite
- WithBottomBar
- Installation
- backend_test.go
- graphify-refresh.sh
- options.go
- KernelDev
- ResolveLibraryClosure
- EventType
- sync.Mutex
- shell
- github.com/stretchr/testify/suite.Suite
- SanitizeText
- taint_victim.c
- changelog.go
- .populateWatchInodes
- Limitations
- Troubleshooting
- IntegrationSuite
- NetworkMonitorRepository
- configedit_test.go
- copyRegularAt
- run
- fileedit_symlink_test.go
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
- openBunTmp
- Execute
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- NewDumpHook
- pathCache
- NewEventFanout
- DaemonEvent
- run
- ReadInstallCert
- SuperblockDev
- Desktop GUI
- IntegrationSuite
- refreshAdmits
- bpfLsmListed
- .eventTypeAllowed
- DiscoverInfraBinaries
- parseResize
- bleedLayout
- TestMain
- LibraryClosure
- Vault
- parseKernelRelease
- runNetworkGuard
- TestCleanRel
- Daemon configuration
- golang.org/x/sys/unix.Timespec
- time.Time
- CleanupStalePins
- .runEditorHarness
- inode_dev.h
- RunWith
- MulticallError
- serve.go

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 109 edges
3. `GuardInodeKey` - 99 edges
4. `Load()` - 89 edges
5. `shell` - 87 edges
6. `absPath()` - 83 edges
7. `IntegrationSuite` - 82 edges
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

## Communities (233 total, 40 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (18): tempJournalRow, Progress, User, invalidPassword, recordTick, launchRule, trustEvent, textAnswer (+10 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (168): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod() (+160 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.17
Nodes (14): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr(), newBoundedKeyFn(), readKey() (+6 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (76): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+68 more)

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
Nodes (25): addErrKind, bpfGuardEvent, BinariesSummary(), classifyAddErr(), decodeGuardEvent(), eventMask(), fsGateLabel(), guardModeKey() (+17 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "BinaryEntry"
Cohesion: 0.08
Nodes (21): deferredBinary, keptEntries, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys(), canonicalBinaryPath(), confinedEntry(), patternMatches() (+13 more)

### Community 11 - "policyEditor"
Cohesion: 0.12
Nodes (13): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), eventsOrAll(), netWarnings() (+5 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (25): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, recordForm, recordReview (+17 more)

### Community 13 - "codeedit.go"
Cohesion: 0.15
Nodes (22): class, pos, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+14 more)

### Community 14 - "CandidateDir"
Cohesion: 0.13
Nodes (13): DiscoverSystem(), expandPlaceholders(), CandidateDir, TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), TestPathsFor(), catalogEntryMatch(), findCatalogRoot() (+5 more)

### Community 15 - "editView"
Cohesion: 0.11
Nodes (8): codeEntry, editView, handButton, Entry, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "runView"
Cohesion: 0.14
Nodes (3): runHost, runView, captureLogs()

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (13): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), ParseMode(), chownUser (+5 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.12
Nodes (56): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+48 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (10): fakeDaemon, Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth_test.go"
Cohesion: 0.15
Nodes (35): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+27 more)

### Community 24 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (73): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+65 more)

### Community 26 - "TrustGuard"
Cohesion: 0.07
Nodes (15): tempGrant, trustHook, trustRows, deleteResKeys(), setSupersedeTrust(), liveTrust(), cStr(), syncMap() (+7 more)

### Community 27 - "Finder"
Cohesion: 0.15
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 28 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (23): recordTickCmd(), runTUI(), recordModel, diffModel, editorModel, Highlighter, syncViewport(), formatAddr() (+15 more)

### Community 29 - "rcUser"
Cohesion: 0.13
Nodes (34): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+26 more)

### Community 30 - "Config"
Cohesion: 0.07
Nodes (36): swapFile(), newControlManager(), startControlManager(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), inspectorPaths(), resolveInspectors() (+28 more)

### Community 31 - "MonitorRepository"
Cohesion: 0.12
Nodes (15): MonitorRepository, NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError() (+7 more)

### Community 32 - "UI"
Cohesion: 0.08
Nodes (23): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestRunWithRefusesConcurrentRun() (+15 more)

### Community 33 - "NetEventType"
Cohesion: 0.07
Nodes (17): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, newNetEventChecks(), FormatAddr(), NetEventTypes(), Ntohs() (+9 more)

### Community 34 - "picker"
Cohesion: 0.06
Nodes (46): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+38 more)

### Community 35 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (13): checkSecret(), newRecordFlow(), File, backend, backendOp, backendState, BinaryCheck, Message (+5 more)

### Community 38 - "controlServer"
Cohesion: 0.17
Nodes (8): readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), controlServer, editControlSession, grantRequest

### Community 39 - "TempBinary"
Cohesion: 0.06
Nodes (27): confirmUserPlaced(), forwardCheck(), tempMaskOp, tempRow, binaryStatOf(), TempBinary, TemporaryGrant, TempRule (+19 more)

### Community 40 - "ConfinedEntry"
Cohesion: 0.13
Nodes (28): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+20 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "runMaintenanceMode"
Cohesion: 0.18
Nodes (19): parseMaintenanceFlags(), runMaintenanceMode(), TestValidateMaintenanceFlags(), validateMaintenanceFlags(), deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup() (+11 more)

### Community 45 - "runDaemon"
Cohesion: 0.07
Nodes (39): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), loadDaemonConfig(), newPinGeneration() (+31 more)

### Community 46 - "User"
Cohesion: 0.13
Nodes (22): entryWriters(), reserveChain(), underResource(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), addKeysToAgent() (+14 more)

### Community 48 - "lockAndDeprovision"
Cohesion: 0.11
Nodes (14): entryOf(), deprovisionKind, migrationTarget, writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr(), classifyMigrationTarget(), isDeprovisionBusy() (+6 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 51 - "cursor.go"
Cohesion: 0.10
Nodes (15): binaryList, buttonCursor, newBinaryList(), handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm() (+7 more)

### Community 52 - "Client"
Cohesion: 0.08
Nodes (16): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Reader (+8 more)

### Community 53 - "watchGroup"
Cohesion: 0.11
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 54 - "tabs.go"
Cohesion: 0.03
Nodes (19): appTheme, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar (+11 more)

### Community 55 - "Plan"
Cohesion: 0.11
Nodes (17): InsertSectionsBeforeLibraries(), TestInsertSectionsBeforeLibraries(), Plan, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan() (+9 more)

### Community 56 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 57 - "mutedText"
Cohesion: 0.08
Nodes (39): binaryRow, homePage, iconTile, multiAnswer, panel, pill, statusTile, targetRow (+31 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (52): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+44 more)

### Community 59 - "os.File"
Cohesion: 0.05
Nodes (58): ReplacementCheck, VettedInode, Multicall, MulticallKind, verifiedBinary(), ExeKey(), IsMulticallRefusal(), multicallFamily() (+50 more)

### Community 60 - "Candidate"
Cohesion: 0.05
Nodes (56): confirmMasterKeyOverwrite(), runGenKey(), addedResources(), appendSectionsAndEdit(), applyDiffAdditions(), collectDiffAdditions(), configPaths(), pathCovered() (+48 more)

### Community 61 - "Recorder"
Cohesion: 0.17
Nodes (15): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), sortedKeys() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+7 more)

### Community 63 - "Network"
Cohesion: 0.12
Nodes (8): netGuardPill(), NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy, Network, verifyHashes()

### Community 64 - "RefreshSection"
Cohesion: 0.07
Nodes (36): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+28 more)

### Community 65 - "pinstate.go"
Cohesion: 0.11
Nodes (32): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive() (+24 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "Binary"
Cohesion: 0.12
Nodes (25): lines(), pickLibs(), appletLine(), beyondEvents(), classifyGroup(), classifyLibs(), finishBinary(), Binary (+17 more)

### Community 69 - "FetchLatest"
Cohesion: 0.11
Nodes (20): installHooks(), releaseBuild(), Run(), staged(), channelOf(), downloadFile(), downloadReleaseFiles(), FetchLatest() (+12 more)

### Community 70 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 72 - "Guard"
Cohesion: 0.07
Nodes (7): exeRow, rootHandle, openRoot(), SharedPinDegraded(), TestNoteRootMountFollowsMountID(), Guard, releaseRootTrust()

### Community 73 - "FileTarget"
Cohesion: 0.08
Nodes (18): filesPage, specFromForm(), guardPill(), specSummary(), EventTypes(), maskFor(), newExeResolver(), Files (+10 more)

### Community 74 - "os.FileMode"
Cohesion: 0.12
Nodes (21): dupFD(), installFile(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow(), TestInstallSSHAgentRefusesSymlinkedUnitDir() (+13 more)

### Community 76 - "catalogRefresher"
Cohesion: 0.18
Nodes (7): drain(), newCatalogRefresher(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 77 - "journal"
Cohesion: 0.14
Nodes (7): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 78 - "deploy"
Cohesion: 0.17
Nodes (19): archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), downloadAndVerify(), fetchGUI(), fetchPair(), guiAssetsFor() (+11 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GlobReservations"
Cohesion: 0.14
Nodes (11): GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, resolveBits(), resolveChildren(), statKey() (+3 more)

### Community 81 - "serveWebSocket"
Cohesion: 0.19
Nodes (10): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop() (+2 more)

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - "githubRelease"
Cohesion: 0.15
Nodes (24): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), tags(), TestReleasesSince(), compareStableVersions(), filterChannel() (+16 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "runBackend"
Cohesion: 0.13
Nodes (15): ParseEventsFlag(), reserveStdout(), runBackend(), LiveModeAvailable(), TestValidateFlagsEditConfig(), runEditProtected(), validateFlags(), TestValidateForwardFlags() (+7 more)

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 87 - "recordFlow"
Cohesion: 0.16
Nodes (18): encryptBlocker(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickBinaries(), pickEncryption(), reviewRecording() (+10 more)

### Community 88 - "liveSession"
Cohesion: 0.17
Nodes (11): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), displayEntries(), eventsLabel(), forwardBinaries(), holdForward() (+3 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "inodeChain"
Cohesion: 0.21
Nodes (16): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+8 more)

### Community 92 - "Monitor"
Cohesion: 0.16
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (10): foreignPinnedLink(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), hook (+2 more)

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
Cohesion: 0.17
Nodes (12): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+4 more)

### Community 100 - "highlight_test.go"
Cohesion: 0.09
Nodes (36): newChangelogModel(), rel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog(), TestNewChangelogModelViewContainsTitleAndNotes() (+28 more)

### Community 101 - "Prereq"
Cohesion: 0.13
Nodes (11): classifySetupError(), classifySupportError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), TestClassifySupportErrorEncryptionNotEnabled(), TestClassifySupportErrorEncryptionNotEnabledF2fs(), TestClassifySupportErrorGeneric() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "installPage"
Cohesion: 0.06
Nodes (24): anyDaemonRunning(), actionGate, binSource, checkGrid, installAction, installedBinary, installPage, latestRelease (+16 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.22
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "openFileVaultWithMasterKey"
Cohesion: 0.18
Nodes (14): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), openFileVaultWithMasterKey(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated() (+6 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 110 - "recoverFileInPlace"
Cohesion: 0.31
Nodes (13): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), recoverFileInPlace(), rejectSidecarSymlink(), stageRecovery(), inode() (+5 more)

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "NetGuardEvent"
Cohesion: 0.06
Nodes (28): snapNetGuard, NetGuardEvent, Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace() (+20 more)

### Community 113 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 115 - "readGrantRequest"
Cohesion: 0.18
Nodes (13): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+5 more)

### Community 116 - "Vault"
Cohesion: 0.10
Nodes (21): DefaultNewFileMode(), Vault, IsOctalMode(), Open(), inodeOf(), mustRead(), mustWrite(), openRootT() (+13 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.36
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.06
Nodes (58): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+50 more)

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
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - ".launch"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "deploy"
Cohesion: 0.06
Nodes (45): installGUI(), buildBinaryIfNeeded(), certifyInstall(), deploy(), installBinaryOnly(), markComplete(), preflightDeployedBinary(), installConfig() (+37 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "dialLiveSession"
Cohesion: 0.06
Nodes (36): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+28 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "GuardEvent"
Cohesion: 0.04
Nodes (19): holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), TestLocalizePaths() (+11 more)

### Community 143 - "batchInodeRows"
Cohesion: 0.19
Nodes (12): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+4 more)

### Community 144 - "runOfflineEdit"
Cohesion: 0.17
Nodes (11): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptStep(), pickDirsToDecrypt(), VerifyEncryptedKeys() (+3 more)

### Community 145 - "Open"
Cohesion: 0.16
Nodes (20): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+12 more)

### Community 146 - "ValidateTargets"
Cohesion: 0.18
Nodes (9): existingSection(), guardedMatch(), overlaps(), refusedPath(), validateTarget(), ValidateTargets(), TestValidateTargets(), Existing (+1 more)

### Community 147 - "install_run.go"
Cohesion: 0.09
Nodes (17): adaptiveGrid, gapVBox, logGrid, logHook, logLine, stageItem, stageState, logRow() (+9 more)

### Community 148 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 149 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 150 - "TrustGlob"
Cohesion: 0.29
Nodes (5): symlinkTargetGlob(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 152 - "filevault_test.go"
Cohesion: 0.27
Nodes (17): looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip(), TestVaultUnlockLockDispatchForRegularFile() (+9 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.07
Nodes (41): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), resolveGuardConfig() (+33 more)

### Community 155 - "secureResources"
Cohesion: 0.24
Nodes (8): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 157 - "WithBottomBar"
Cohesion: 0.24
Nodes (7): decryptDirectories(), TestDecryptDirectoriesEmpty(), TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "options.go"
Cohesion: 0.31
Nodes (9): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+1 more)

### Community 162 - "KernelDev"
Cohesion: 0.09
Nodes (30): BinaryStat, isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs(), binaryStatOf() (+22 more)

### Community 163 - "ResolveLibraryClosure"
Cohesion: 0.18
Nodes (11): hostView(), mountVouchesOwnership(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch() (+3 more)

### Community 164 - "EventType"
Cohesion: 0.08
Nodes (19): systemPatterns(), systemRule(), snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents() (+11 more)

### Community 165 - "sync.Mutex"
Cohesion: 0.05
Nodes (47): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+39 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (43): fakeLink, pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines (+35 more)

### Community 167 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.22
Nodes (4): Classify(), TestClassifyVerifiesTheInodeThatRan(), harnessSuite, editorHarness

### Community 168 - "SanitizeText"
Cohesion: 0.04
Nodes (34): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), formatElapsed() (+26 more)

### Community 170 - "changelog.go"
Cohesion: 0.13
Nodes (21): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+13 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "NetworkMonitorRepository"
Cohesion: 0.20
Nodes (3): NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase()

### Community 176 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 177 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 178 - "run"
Cohesion: 0.22
Nodes (10): confirmKeyDeletion(), run(), runUninstall(), RunWith(), removeKeyAndEmptyDir(), removeMasterKey(), revertBunLaunchers(), TestRemoveKeyAndEmptyDir() (+2 more)

### Community 179 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.09
Nodes (27): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), MakeDisplayPaths(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe() (+19 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "openBunTmp"
Cohesion: 0.33
Nodes (9): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp (+1 more)

### Community 195 - "Execute"
Cohesion: 0.25
Nodes (7): Requested(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger(), subcommand(), main()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `--record` (+2 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 202 - "DaemonEvent"
Cohesion: 0.09
Nodes (36): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+28 more)

### Community 203 - "run"
Cohesion: 0.36
Nodes (8): InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestReadInstalledVersion(), GUIInstalled()

### Community 204 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 205 - "SuperblockDev"
Cohesion: 0.29
Nodes (7): FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev(), SuperblockDevPath(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "refreshAdmits"
Cohesion: 0.33
Nodes (6): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), BinaryRule, materialized(), Extension

### Community 210 - "bpfLsmListed"
Cohesion: 0.40
Nodes (6): bpfLsmListed(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 213 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 215 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 216 - "LibraryClosure"
Cohesion: 0.10
Nodes (19): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+11 more)

### Community 218 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 219 - "runNetworkGuard"
Cohesion: 0.10
Nodes (20): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI() (+12 more)

### Community 221 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 224 - "time.Time"
Cohesion: 0.11
Nodes (16): LockedError, addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+8 more)

### Community 227 - "CleanupStalePins"
Cohesion: 0.40
Nodes (5): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 229 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 230 - "RunWith"
Cohesion: 0.14
Nodes (17): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr() (+9 more)

### Community 237 - "serve.go"
Cohesion: 0.15
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

## Knowledge Gaps
- **31 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+26 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 588 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **40 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `liveSession`, `runMaintenanceMode`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _31 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08959842325695984 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.047) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0165437512457644 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `runMaintenanceMode` to `daemon/daemon.go`, `edit-protected`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._