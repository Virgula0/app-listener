# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 436 files · ~1,614,510 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 5513 nodes · 23421 edges · 199 communities (169 shown, 30 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2052 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fe1504aa`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- RefreshSection
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- catalogrefresh.go
- planUpdaters
- ui.go
- sanitizeTerminalText
- system.go
- orphans_test.go
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- Hash
- NetGuardEvent
- guard.bpf.c
- SanitizeText
- networkPage
- inodeChain
- codeedit.go
- newTestTarget
- Serve
- bpfstats/main.go
- render.go
- icon
- github.com/spf13/cobra.Command
- IntegrationSuite
- Guard
- controlServer
- fileedit_model_test.go
- .objs
- IntegrationSuite
- absPath
- monitorUnitTest
- IntegrationSuite
- update.go
- vaultFS
- stdio
- golang.org/x/sys/unix.Stat_t
- Config
- LibraryClosure
- fakeGuardRepo
- rootHandle
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- install/install.go
- uninstall_system.go
- GuardInodeKey
- NewGuard
- buildOneGuard
- .Refresh
- exe_supersede.h
- buildTrustedSet
- runView
- pinstate.go
- gui/session.go
- runMonitor
- changelog.go
- FileTarget
- fcntl
- copytree.go
- fileedit_symlink_test.go
- time.Time
- btrfsLayoutFrom
- Network
- sync.Mutex
- inode_dev
- NetEventType
- process_vm_readv.c
- GuardEvent
- .ExpandWhitelist
- model
- wizard.go
- Monitor
- runNetworkGuard
- github.com/charmbracelet/bubbles/textarea.Model
- liveSession
- highlight_test.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- netModel
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- headless/headless_test.go
- policyEditor
- check_watched
- Well-known bypass classes
- daemonModel
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- github.com/charmbracelet/bubbletea.KeyMsg
- Build & sign release assets (reusable workflow)
- installPage
- KernelDev
- mountTable
- NetworkGuardRepository
- app-listener Terminal Demo Recording
- IntegrationSuite
- IntegrationSuite
- GenerateConf
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- bufio.Reader
- New
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- CandidateDir
- make test-integration (rootful Docker suite)
- IntegrationSuite
- supersede_probe.c
- .createEntry
- trace-app-libs.sh
- IntegrationSuite
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- EventType
- io.Reader
- tempGrantSetup
- pathCache
- NewEventFanout
- net_tester/main.go
- .eventTypeAllowed
- newTestVetter
- Vault
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- FileEvent
- layout.go
- options.go
- IntegrationSuite
- Installation
- os.File
- graphify-refresh.sh
- time.Duration
- parseGuardEvent
- TestMain
- WithSuspectsFrom
- feed.go
- shell
- unixModeToFileMode
- CleanupStalePins
- inode_dev.h
- RefreshCatalog
- IntegrationSuite
- Limitations
- Troubleshooting
- DiscoverInfraBinaries
- mountsUnitTest
- Daemon configuration
- edit-protected
- parseResize
- configedit_test.go
- Compatibility
- How it works
- writeWithin
- launchKey
- bleedLayout
- launch_rule_at
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .runEditorHarness
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- runDaemon
- monitor: see what touches a path
- kernel-prof.sh
- .populateWatchInodes
- Desktop GUI

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 96 edges
3. `GuardInodeKey` - 96 edges
4. `Load()` - 90 edges
5. `IntegrationSuite` - 82 edges
6. `absPath()` - 81 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 68 edges
9. `shell` - 65 edges
10. `BinaryEntry` - 60 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (199 total, 30 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (21): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, textAnswer, btrfsIoctlFsInfoArgs, multicallCacheKey (+13 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (102): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered(), TestIsEncryptedRegularFile(), TestFscryptHarness() (+94 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "RefreshSection"
Cohesion: 0.07
Nodes (38): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, findSectionEnd(), findSectionStart(), firstContentLine(), isLibDirectiveText() (+30 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+70 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.11
Nodes (34): code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image(), fill_inode_key(), get_task_exe_key() (+26 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.07
Nodes (66): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir() (+58 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (17): addErrKind, bpfGuardEvent, BinariesSummary(), classifyAddErr(), commMatchesGuardedBinary(), eventMask(), guardModeKey(), addRecorder() (+9 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "catalogrefresh.go"
Cohesion: 0.31
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "ui.go"
Cohesion: 0.07
Nodes (34): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), needsBinary() (+26 more)

### Community 13 - "sanitizeTerminalText"
Cohesion: 0.17
Nodes (8): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), TestSanitizeTerminalText()

### Community 14 - "system.go"
Cohesion: 0.07
Nodes (41): TestDaemonUnitWithMetadataOutput(), dupFD(), installFile(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit() (+33 more)

### Community 15 - "orphans_test.go"
Cohesion: 0.18
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 19 - "fileEditModel"
Cohesion: 0.15
Nodes (6): clipLabel(), humanSize(), chownUser, fileEditMode, fileEditModel, fileKind

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (12): ownerPaths(), fakeDaemon, Resource, GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase (+4 more)

### Community 23 - "Hash"
Cohesion: 0.13
Nodes (28): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+20 more)

### Community 24 - "NetGuardEvent"
Cohesion: 0.18
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (53): add_inode_to_guard(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class(), evict_inode_from_guard() (+45 more)

### Community 26 - "SanitizeText"
Cohesion: 0.05
Nodes (36): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), candidates(), confirmable() (+28 more)

### Community 27 - "networkPage"
Cohesion: 0.18
Nodes (6): networkPage, watchTab, entryPaths(), newNetGuardTab(), optionHint(), netGuardPill()

### Community 28 - "inodeChain"
Cohesion: 0.26
Nodes (13): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+5 more)

### Community 29 - "codeedit.go"
Cohesion: 0.16
Nodes (27): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), moveTo(), Navigate() (+19 more)

### Community 30 - "newTestTarget"
Cohesion: 0.09
Nodes (22): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+14 more)

### Community 31 - "Serve"
Cohesion: 0.10
Nodes (26): basicAuth(), configureWebSocket(), isInteractiveTerminal(), newServeListener(), newServeServer(), readResizeLoop(), relayServeSignals(), runServe() (+18 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - "render.go"
Cohesion: 0.14
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 34 - "icon"
Cohesion: 0.11
Nodes (12): binaryList, filesPage, newBinaryList(), newBinaryRow(), specFromForm(), newTargetView(), specSummary(), asset() (+4 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.07
Nodes (40): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+32 more)

### Community 37 - "Guard"
Cohesion: 0.05
Nodes (41): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, deferredBinary, HashStore, sharedHashEntry (+33 more)

### Community 38 - "controlServer"
Cohesion: 0.18
Nodes (8): readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), controlServer, editControlSession, grantRequest

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.27
Nodes (19): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+11 more)

### Community 40 - ".objs"
Cohesion: 0.10
Nodes (7): exeRow, tempGrant, tempMaskOp, tempRow, deleteResKeys(), supersededMark(), liveTrust()

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (5): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "monitorUnitTest"
Cohesion: 0.12
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 44 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 45 - "update.go"
Cohesion: 0.06
Nodes (53): applyUpdate(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), confirmUpdate(), downloadAndVerify(), downloadFile(), downloadReleaseFiles() (+45 more)

### Community 46 - "vaultFS"
Cohesion: 0.13
Nodes (13): TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), isBinaryContent(), listEntries(), openRootT(), TestDefaultNewFileMode(), TestIsBinaryContent() (+5 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 49 - "Config"
Cohesion: 0.06
Nodes (48): drain(), newCatalogRefresher(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh() (+40 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 52 - "rootHandle"
Cohesion: 0.17
Nodes (5): rootHandle, openRoot(), TestNoteRootMountFollowsMountID(), TestReadOnlyWalkToleratesDanglingSymlinks(), PathUnder()

### Community 53 - "daemonconfig.go"
Cohesion: 0.14
Nodes (36): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+28 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (15): diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), assertViewFits() (+7 more)

### Community 55 - "TrustGuard"
Cohesion: 0.08
Nodes (17): trustHook, trustRows, dropSupersededLocked(), PruneSuperseded(), setSupersedeTrust(), supersedeMaps(), supersedePruneLoop(), warnLostStamps() (+9 more)

### Community 56 - "install/install.go"
Cohesion: 0.10
Nodes (28): runDiffCatalog(), applyLiveRefresh(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), encryptDirectories(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256() (+20 more)

### Community 57 - "uninstall_system.go"
Cohesion: 0.04
Nodes (75): reloadDaemonForPasswordChange(), buildBinaryIfNeeded(), installBinaryOnly(), preflightDeployedBinary(), confirmRunPrereq(), resolveFilesystemPrereqs(), installConfig(), mustCwd() (+67 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (49): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+41 more)

### Community 59 - "NewGuard"
Cohesion: 0.16
Nodes (12): resolveGuardConfig(), runGuard(), runGuardTUI(), TestWithChmodDropWriteNeedsReadOnly(), canonicalPaths(), Mode, modeString(), NewGuard() (+4 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.09
Nodes (38): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+30 more)

### Community 61 - ".Refresh"
Cohesion: 0.04
Nodes (11): appTheme, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, newNavItem() (+3 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 63 - "buildTrustedSet"
Cohesion: 0.14
Nodes (17): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+9 more)

### Community 64 - "runView"
Cohesion: 0.07
Nodes (17): installUI, logHook, logLine, runView, stageItem, stageState, scrollable(), logImportance() (+9 more)

### Community 65 - "pinstate.go"
Cohesion: 0.06
Nodes (53): CheckBPFLSM(), configureDaemonLogging(), confirmMasterKeyOverwrite(), loadDaemonConfig(), lockOneRoot(), lockRootRecovering(), newPinGeneration(), prepareDaemonStart() (+45 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (40): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 67 - "runMonitor"
Cohesion: 0.08
Nodes (24): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), pathWithinDir(), ResolveTargets(), validatePair() (+16 more)

### Community 68 - "changelog.go"
Cohesion: 0.14
Nodes (21): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), showChangelog(), rel(), tags() (+13 more)

### Community 69 - "FileTarget"
Cohesion: 0.14
Nodes (10): Target, guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec() (+2 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copytree.go"
Cohesion: 0.11
Nodes (24): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), copyDirContents(), copyEntry(), copyRegularAt() (+16 more)

### Community 72 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 73 - "time.Time"
Cohesion: 0.09
Nodes (22): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+14 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "Network"
Cohesion: 0.14
Nodes (9): newExeResolver(), guardRow(), clonePolicy(), NetExtra, NetPolicy, Network, NewNetwork(), exeCached (+1 more)

### Community 76 - "sync.Mutex"
Cohesion: 0.17
Nodes (5): swapFile(), newControlManager(), startControlManager(), configEditor, controlManager

### Community 77 - "inode_dev"
Cohesion: 0.18
Nodes (15): exe_is_superseded(), exe_refused(), inode_dev(), check_and_emit_args(), code_suspect(), current_is_inspector(), exe_tag_bits(), fill_path() (+7 more)

### Community 78 - "NetEventType"
Cohesion: 0.06
Nodes (19): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), Cstr(), FormatAddr() (+11 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.09
Nodes (15): holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), GuardEvent (+7 more)

### Community 82 - "model"
Cohesion: 0.22
Nodes (7): formatType(), listenForEvents(), syncViewport(), tickStats(), eventLine, eventMsg, model

### Community 83 - "wizard.go"
Cohesion: 0.07
Nodes (39): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), deletePostBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories() (+31 more)

### Community 84 - "Monitor"
Cohesion: 0.16
Nodes (4): NewMonitor(), RealFileEngines(), Monitor, watchedInode

### Community 85 - "runNetworkGuard"
Cohesion: 0.12
Nodes (20): CheckEBPF(), ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard() (+12 more)

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.14
Nodes (19): findButton, span, column(), findAll(), Finder, lowerRunes(), NewFinder(), overlay() (+11 more)

### Community 87 - "liveSession"
Cohesion: 0.12
Nodes (21): auditAfterEdit(), dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), displayEntries() (+13 more)

### Community 88 - "highlight_test.go"
Cohesion: 0.15
Nodes (17): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+9 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 92 - "DaemonEvent"
Cohesion: 0.12
Nodes (31): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+23 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (12): reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+4 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "headless/headless_test.go"
Cohesion: 0.40
Nodes (3): Check(), buildTiny(), TestCheck()

### Community 100 - "policyEditor"
Cohesion: 0.12
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 101 - "check_watched"
Cohesion: 0.24
Nodes (14): admit_refusal(), check_watched(), emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg() (+6 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "daemonModel"
Cohesion: 0.20
Nodes (11): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+3 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.05
Nodes (39): binaryRow, binSource, checkGrid, handButton, homePage, iconTile, installedBinary, installPage (+31 more)

### Community 110 - "KernelDev"
Cohesion: 0.11
Nodes (27): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), KernelDev(), StatBinary(), StatBinaryFile() (+19 more)

### Community 111 - "mountTable"
Cohesion: 0.14
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), ReadMountinfo(), unescapeMountinfo(), mountTable (+1 more)

### Community 112 - "NetworkGuardRepository"
Cohesion: 0.12
Nodes (15): NetworkGuardRepository, NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo() (+7 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 116 - "GenerateConf"
Cohesion: 0.12
Nodes (27): TestDiffMergePreservesExistingAndParses(), askEncryption(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestGeneratedLibraryBlockRoundTrips() (+19 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "bufio.Reader"
Cohesion: 0.20
Nodes (11): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest(), TestReadGrantRequestForward() (+3 more)

### Community 121 - "New"
Cohesion: 0.07
Nodes (49): collectFilesystemPrereqs(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), harnessSuite, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder() (+41 more)

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

### Community 127 - "CandidateDir"
Cohesion: 0.06
Nodes (34): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, GlobChild, statFunc, compileGlobNames() (+26 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "EventType"
Cohesion: 0.07
Nodes (27): EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), binaryStatOf(), TempBinary, TemporaryGrant, TempRule (+19 more)

### Community 143 - "io.Reader"
Cohesion: 0.50
Nodes (3): BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype()

### Community 144 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 145 - "pathCache"
Cohesion: 0.60
Nodes (3): evalSymlinksOrEmpty(), pathCache, pathCacheEntry

### Community 146 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 149 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 150 - "Vault"
Cohesion: 0.12
Nodes (18): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+10 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.24
Nodes (21): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+13 more)

### Community 154 - "FileEvent"
Cohesion: 0.11
Nodes (7): snapEngines, snapMonitor, BpfEvent, eventUnitTest, FileEvent, DecodeBpfEvent(), fakeMonitorRepo

### Community 155 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 156 - "options.go"
Cohesion: 0.31
Nodes (6): installAction, InstallMode, InstallRequest, Options, InstallHooks, Run()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.04
Nodes (74): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+66 more)

### Community 161 - "time.Duration"
Cohesion: 0.14
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 162 - "parseGuardEvent"
Cohesion: 0.29
Nodes (8): decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels(), suspectLabel()

### Community 163 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 164 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 165 - "feed.go"
Cohesion: 0.07
Nodes (43): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+35 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (26): pageDef, pageID, resizeWatch, shell, sizedDialog, sysStatus, newLogo(), newShell() (+18 more)

### Community 167 - "unixModeToFileMode"
Cohesion: 0.29
Nodes (6): fileModeToUnixMode(), isOctalMode(), TestFileModeToUnixModeRoundTrip(), TestIsOctalMode(), TestUnixModeToFileMode(), unixModeToFileMode()

### Community 168 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 169 - "inode_dev.h"
Cohesion: 0.24
Nodes (8): ctx_ptr(), sb_dev(), btrfs_copy_gate(), guard_file_ioctl(), guard_file_ioctl_compat(), guard_inode_free(), is_btrfs_copy_ioctl(), trust_inode_free()

### Community 170 - "RefreshCatalog"
Cohesion: 0.17
Nodes (14): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), RefreshCatalog(), staleDiscordLibConf() (+6 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 179 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 185 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

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

### Community 196 - "runDaemon"
Cohesion: 0.18
Nodes (12): awaitStartupOrSignal(), catchLifecycleSignals(), markReady(), notifySystemdReady(), runDaemon(), startGuardedDaemonAbortable(), startPprof(), writePidFile() (+4 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 206 - "Desktop GUI"
Cohesion: 0.22
Nodes (9): Actions, Desktop GUI, Files, Home, Install, Network, Tables, The daemon binary (+1 more)

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 521 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **30 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `time.Time` to `Daemon configuration`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10559370136409411 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `time.Time`, `README.md`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.026979368718039146 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `time.Time` to `globreserve_test.go`, `catalogrefresh.go`, `Config`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._