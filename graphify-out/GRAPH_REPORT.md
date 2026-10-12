# Graph Report - app-listener  (2026-10-12)

## Corpus Check
- 365 files · ~1,561,399 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4686 nodes · 19045 edges · 202 communities (169 shown, 33 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1361 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d5c710d8`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- orphans_test.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- mkdirs
- runGuard
- planUpdaters
- runNetworkGuard
- Candidate
- install.go
- filevault.go
- .startGuardStd
- guard_trust.bpf.c
- confgen.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- buildOneGuard
- guard.bpf.c
- NetworkMonitor
- find_test.go
- Ledger
- sanitizeTerminalText
- Monitor
- serve.go
- GuardInodeKey
- IntegrationSuite
- fileedit_model_test.go
- IntegrationSuite
- dataRowRenderer
- uninstall_system.go
- controlServer
- guardSpec
- Finder
- New
- IntegrationSuite
- liveSession
- IntegrationSuite
- update.go
- highlight_test.go
- .syncSetsLocked
- TrustGuard
- tempRow
- inode_dev
- IntegrationSuite
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- absPath
- codeedit.go
- sync.Mutex
- engine
- GuardEvent
- auditGroupCoverage
- headerWidget
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- catalogrefresh.go
- globBuilder
- Config
- time.Duration
- fcntl
- copyRegularAt
- integrationtests/networkguard_test.go
- watchSet
- btrfsLayoutFrom
- monitorUnitTest
- github.com/spf13/cobra.Command
- IntegrationSuite
- buildTrustedSet
- process_vm_readv.c
- golang.org/x/sys/unix.Stat_t
- SuperblockDev
- model
- newTestVetter
- MountVouch
- prepareDaemonStart
- github.com/charmbracelet/bubbles/textarea.Model
- fakeDaemon
- mustSubkey
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- safeio.go
- KernelDev
- IntegrationSuite
- Well-known bypass classes
- runNetworkMonitor
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- expandPlaceholders
- Build & sign release assets (reusable workflow)
- runMonitor
- fileedit_symlink_test.go
- net_tester/main.go
- NewEventFanout
- app-listener Terminal Demo Recording
- IntegrationSuite
- tempGrantSetup
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- catalogRefresher
- NewDumpHook
- WebSocket /ws client
- netModel
- string
- install.sh
- mc_tag.h
- changelog.go
- make test-integration (rootful Docker suite)
- EventType
- supersede_probe.c
- bufio.Reader
- trace-app-libs.sh
- IntegrationSuite
- serveWebSocket
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- CandidateDir
- SanitizeText
- preload_wait.c
- Guard
- writeWithin
- runUpdate
- Run
- uninstall.go
- StripPinnedTempAllows
- daemon: permanent protection with encryption at rest
- BinaryEntry
- testKeyPair
- inodeChain
- New
- TestMonitorUseCaseLifecycle
- Installation
- os.File
- graphify-refresh.sh
- bpfstats/main.go
- ResolveLibraryClosure
- find.go
- time.Time
- parseGuardEvent
- releaseAssets
- bottomBar
- IntegrationSuite
- BinaryRule
- fakeVault
- Limitations
- Troubleshooting
- stripTempJournal
- .runAttrOp
- Daemon configuration
- edit-protected
- IntegrationSuite
- lockRootRecovering
- Compatibility
- How it works
- findTopLevelWindow
- markReady
- bpfLsmListed
- downloadFile
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- vouchedDevs
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- DiscoverInfraBinaries
- formatEntry
- TestBatchInodeRows
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS
- Vault
- newTestSession

## God Nodes (most connected - your core abstractions)
1. `Config` - 95 edges
2. `Guard` - 95 edges
3. `GuardInodeKey` - 93 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 66 edges
9. `fileEditModel` - 60 edges
10. `User` - 56 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `Well-known bypass classes` --references--> `btrfs_search POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `Well-known bypass classes` --references--> `io_uring POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared bubbletea TUI layout across monitor and guard modes** — media_demo_monitor_tui, media_demo_guard_whitelist_tui, media_demo_guard_blacklist_tui, media_demo_tui_resource_footer [INFERRED 0.85]
- **Reserved-name trust defenses (trust object #3)** — integrationtests_exploits_readme_glob_plant, integrationtests_exploits_readme_lib_probe, readme_bun_tmpdir_redirect, claude_catalog_refresh [INFERRED 0.85]
- **Release build, sign and verify pipeline** — _github_workflows_release_prerelease_workflow, _github_workflows_promote_release_promote_release_workflow, _github_workflows__build_release_build_release_workflow, _github_workflows__build_release_ed25519_signing, readme_install_sh, readme_update [INFERRED 0.95]
- **Verifier budget measurement and gating** — _github_workflows_verifier_verifier_gate, claude_bpfstats, readme_daemon_check, claude_verifier_complexity_budget, _github_workflows__build_release_ebpf_regeneration [INFERRED 0.95]

## Communities (202 total, 33 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (22): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+14 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (110): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+102 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.05
Nodes (14): IntegrationSuite, IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix() (+6 more)

### Community 3 - "orphans_test.go"
Cohesion: 0.18
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (82): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+74 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.09
Nodes (53): fileID(), trustManager, openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), refreshAdmits() (+45 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (16): addErrKind, guardUnitTest, classifyAddErr(), eventMask(), guardModeKey(), addRecorder(), copySelf(), runTool() (+8 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "runGuard"
Cohesion: 0.16
Nodes (13): selfProtectSpecs(), displayEntries(), holdForward(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless() (+5 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "runNetworkGuard"
Cohesion: 0.10
Nodes (20): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), runHeadless() (+12 more)

### Community 13 - "Candidate"
Cohesion: 0.10
Nodes (30): appendSectionsAndEdit(), collectDiffAdditions(), pathCovered(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), cleanOrphanedFscrypt() (+22 more)

### Community 14 - "install.go"
Cohesion: 0.08
Nodes (40): reloadDaemonForPasswordChange(), setupBunTmpdir(), promptEditPassword(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), deploy(), ensureInstalledBinary() (+32 more)

### Community 15 - "filevault.go"
Cohesion: 0.05
Nodes (55): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), resolveFilesystemPrereqs(), askEncryption(), lockVaultFully(), unlockUnderGuard() (+47 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "confgen.go"
Cohesion: 0.07
Nodes (52): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, Section, ConfSafePath(), findSectionEnd(), findSectionStart() (+44 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.25
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (43): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+35 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (12): ownerPaths(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles() (+4 more)

### Community 23 - "auth.go"
Cohesion: 0.07
Nodes (44): ParseEventsFlag(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+36 more)

### Community 24 - "buildOneGuard"
Cohesion: 0.09
Nodes (35): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded() (+27 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (59): add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial() (+51 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.09
Nodes (14): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), NewNetworkMonitor() (+6 more)

### Community 27 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 28 - "Ledger"
Cohesion: 0.12
Nodes (10): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), Ledger, Pending (+2 more)

### Community 29 - "sanitizeTerminalText"
Cohesion: 0.17
Nodes (8): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), TestSanitizeTerminalText()

### Community 30 - "Monitor"
Cohesion: 0.12
Nodes (5): evalSymlinksOrEmpty(), NewMonitor(), Monitor, pathCache, pathCacheEntry

### Community 31 - "serve.go"
Cohesion: 0.12
Nodes (22): basicAuth(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop(), relayServeSignals() (+14 more)

### Community 32 - "GuardInodeKey"
Cohesion: 0.06
Nodes (29): refusedReplacements, appendInodeRow(), batchInodeRows(), inodeRowsOf(), engine, Inspector, SetInspectors(), engine (+21 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.11
Nodes (35): AtomicWriteAt(), OpenRegularNoFollow(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf() (+27 more)

### Community 35 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 37 - "uninstall_system.go"
Cohesion: 0.09
Nodes (36): deletePostBackups(), restoreBackups(), pickUsers(), offerBackupCleanup(), runUninstall(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir() (+28 more)

### Community 38 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 39 - "guardSpec"
Cohesion: 0.21
Nodes (12): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+4 more)

### Community 41 - "New"
Cohesion: 0.11
Nodes (27): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips() (+19 more)

### Community 45 - "update.go"
Cohesion: 0.13
Nodes (10): compareStableVersions(), newerThanStable(), parseStableVersion(), parseVerPart(), readInstalledVersion(), TestFilterChannel(), TestNewerThanStable(), TestParseStableVersion() (+2 more)

### Community 46 - "highlight_test.go"
Cohesion: 0.14
Nodes (17): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+9 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "TrustGuard"
Cohesion: 0.07
Nodes (16): trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins() (+8 more)

### Community 50 - "tempRow"
Cohesion: 0.18
Nodes (6): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust()

### Community 51 - "inode_dev"
Cohesion: 0.17
Nodes (15): exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), current_is_inspector(), exe_tag_bits(), get_current_exe_inode(), get_task_exe_inode() (+7 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.15
Nodes (36): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+28 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.11
Nodes (7): clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 55 - "absPath"
Cohesion: 0.08
Nodes (5): IntegrationSuite, absPath(), TestMain(), helperChild(), TestMain()

### Community 56 - "codeedit.go"
Cohesion: 0.25
Nodes (17): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+9 more)

### Community 57 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 58 - "engine"
Cohesion: 0.08
Nodes (11): TestNearestRootClaims(), TestTieRank(), ConfigureSharedPinning(), engine, nearestRootClaims(), pathWithin(), TestPathWithinMatchesPrefixRule(), tieRank() (+3 more)

### Community 59 - "GuardEvent"
Cohesion: 0.11
Nodes (15): GuardEvent, logBacklogDenial(), Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace() (+7 more)

### Community 60 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 61 - "headerWidget"
Cohesion: 0.18
Nodes (6): columnState, dataRow, headerWidget, newColumnState(), newDataRow(), newHeaderWidget()

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 64 - "errno"
Cohesion: 0.28
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.27
Nodes (15): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+7 more)

### Community 67 - "globBuilder"
Cohesion: 0.20
Nodes (9): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), inOwnTree() (+1 more)

### Community 68 - "Config"
Cohesion: 0.16
Nodes (20): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager, startTrustGuard(), trustStartupError(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor() (+12 more)

### Community 69 - "time.Duration"
Cohesion: 0.15
Nodes (8): debouncer, New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "integrationtests/networkguard_test.go"
Cohesion: 0.27
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 73 - "watchSet"
Cohesion: 0.20
Nodes (9): addSystemBinary(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern, watchPos (+1 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "monitorUnitTest"
Cohesion: 0.10
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.17
Nodes (16): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+8 more)

### Community 78 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.17
Nodes (12): placedWalk, binaryStatOf(), confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd() (+4 more)

### Community 81 - "SuperblockDev"
Cohesion: 0.17
Nodes (11): checkBtrfsKeys(), BtrfsMounted(), FirstOnBtrfs(), mountinfoDev(), mountinfoHasFstype(), OnBtrfs(), SuperblockDev(), SuperblockDevPath() (+3 more)

### Community 82 - "model"
Cohesion: 0.22
Nodes (7): formatType(), listenForEvents(), syncViewport(), tickStats(), eventLine, eventMsg, model

### Community 83 - "newTestVetter"
Cohesion: 0.18
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 84 - "MountVouch"
Cohesion: 0.07
Nodes (19): GlobChild, statFunc, deleteResKeys(), syncU32Map(), compileGlobNames(), globText(), TrustGuard, resolveBits() (+11 more)

### Community 85 - "prepareDaemonStart"
Cohesion: 0.14
Nodes (18): CheckBPFLSM(), configureDaemonLogging(), confirmMasterKeyOverwrite(), loadDaemonConfig(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck() (+10 more)

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (13): lineStyles, viewportState, baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 88 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.10
Nodes (19): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+11 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.08
Nodes (41): UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemon(), runDaemonUI(), runHeadless() (+33 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.13
Nodes (7): eventsetSummary(), eventTypeKey(), NetGuard, mcNameKey(), modeLabel(), NewNetGuard(), statInodeKey()

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "safeio.go"
Cohesion: 0.07
Nodes (39): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+31 more)

### Community 100 - "KernelDev"
Cohesion: 0.05
Nodes (44): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, fdStat, BinaryStat, DurableStat (+36 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "runNetworkMonitor"
Cohesion: 0.14
Nodes (9): CheckEBPF(), ParseNetEventsFlag(), runNetworkMonitor(), runTUI(), PrintLogo(), NetworkMonitorRepository, NewNetModel(), NetworkMonitorUseCase (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "expandPlaceholders"
Cohesion: 0.21
Nodes (7): expandPlaceholders(), CandidateDir, CandidateDir, LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.10
Nodes (19): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+11 more)

### Community 110 - "fileedit_symlink_test.go"
Cohesion: 0.23
Nodes (11): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+3 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NewEventFanout"
Cohesion: 0.22
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "catalogRefresher"
Cohesion: 0.18
Nodes (6): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), catalogRefresher

### Community 121 - "NewDumpHook"
Cohesion: 0.16
Nodes (15): TestSetupDumpLogRefusesSymlink(), applyVerbosity(), effectiveVerbose(), setupDumpLog(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog() (+7 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 124 - "string"
Cohesion: 0.16
Nodes (6): main(), report(), copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "changelog.go"
Cohesion: 0.18
Nodes (11): changelogText(), newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes() (+3 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "EventType"
Cohesion: 0.10
Nodes (13): BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock(), Run() (+5 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "serveWebSocket"
Cohesion: 0.21
Nodes (7): configureWebSocket(), sameOrigin(), serveWebSocket(), TestSameOrigin(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.13
Nodes (16): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess, logTempGrant() (+8 more)

### Community 143 - "CandidateDir"
Cohesion: 0.21
Nodes (8): confSafeMatch(), BinaryRule, CandidateDir, homeMatchConfined(), symlinkStaysInParent(), boundedGlob(), exists(), matchDir()

### Community 144 - "SanitizeText"
Cohesion: 0.14
Nodes (10): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), binaryVetter, vetDecision, SetHashStore(), TrustGuard, launchKey() (+2 more)

### Community 146 - "Guard"
Cohesion: 0.06
Nodes (17): deferredBinary, rootHandle, openRoot(), TestWithChmodDropWriteNeedsReadOnly(), SharedPinDegraded(), canonicalBinaryPath(), confinedEntry(), Guard (+9 more)

### Community 147 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 148 - "runUpdate"
Cohesion: 0.24
Nodes (13): downloadAndVerify(), fetchReleases(), filterChannel(), latestTime(), newerThanInstalled(), pickLatestRelease(), runUpdate(), TestFetchReleases() (+5 more)

### Community 150 - "uninstall.go"
Cohesion: 0.05
Nodes (64): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), ensureMasterKey(), confirmRunPrereq(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories() (+56 more)

### Community 151 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "BinaryEntry"
Cohesion: 0.09
Nodes (34): binaryVerifyState, HashStore, sharedHashEntry, BinaryEntry, BinariesSummary(), canonicalPaths(), commMatchesGuardedBinary(), ConfinedEntry() (+26 more)

### Community 154 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 155 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 156 - "New"
Cohesion: 0.17
Nodes (13): diffModel, editorModel, moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys() (+5 more)

### Community 157 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.33
Nodes (8): NewGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkMonitorUseCaseLifecycle(), fakeNetworkGuardRepo

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.05
Nodes (60): checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes(), ReplacementCheck, VettedInode, Multicall, MulticallKind (+52 more)

### Community 161 - "bpfstats/main.go"
Cohesion: 0.38
Nodes (9): chownToSudoUser(), compare(), dumpFullLog(), dumpTrace(), loadOne(), main(), measure(), resizeMaps() (+1 more)

### Community 162 - "ResolveLibraryClosure"
Cohesion: 0.18
Nodes (11): hostView(), mountVouchesOwnership(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch() (+3 more)

### Community 163 - "find.go"
Cohesion: 0.20
Nodes (7): findButton, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 164 - "time.Time"
Cohesion: 0.27
Nodes (6): newWatchSet(), fakeInotify, sweepID, sweepMemo, Guard, dirEntry

### Community 165 - "parseGuardEvent"
Cohesion: 0.22
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 166 - "releaseAssets"
Cohesion: 0.22
Nodes (10): assetsFor(), checkUpdatePreconditions(), resolveAssets(), sanityCheckBinary(), TestAssetsFor(), TestResolveAssetsPicksArch(), TestSanityCheckBinaryRefusesOtherMachine(), TestValidateChannel() (+2 more)

### Community 167 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 169 - "BinaryRule"
Cohesion: 0.25
Nodes (9): systemPatterns(), systemRule(), grantNewBinaries(), Line(), Lines(), TestLinesKeyOnLinkAndIncludePending(), BinaryRule, parseEvents() (+1 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "stripTempJournal"
Cohesion: 0.32
Nodes (8): ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders(), ensureHashFilePlaceholder(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace(), writeTempJournal()

### Community 175 - ".runAttrOp"
Cohesion: 0.25
Nodes (4): eventsForPath(), guardEventTypes(), parseGuardEvents(), guardEvent

### Community 176 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "lockRootRecovering"
Cohesion: 0.29
Nodes (8): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), PinPrefix(), sanitizeGen(), SharedPinPrefix(), TestPinPrefixDistinctAndStable(), TestPinPrefixShape()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "findTopLevelWindow"
Cohesion: 0.38
Nodes (4): findTopLevelWindow(), floatWindow(), internAtom(), setFloatHint()

### Community 183 - "markReady"
Cohesion: 0.22
Nodes (8): markReady(), noteReloadedSettings(), notifySystemdReady(), engine, ConfigureInodeCapacity(), InodeOccupancy(), TestConfigureInodeCapacityBounds(), WarnInodeOccupancy()

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 185 - "downloadFile"
Cohesion: 0.32
Nodes (5): downloadFile(), downloadReleaseFiles(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), progressReader

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

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

### Community 194 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 195 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 196 - "TestBatchInodeRows"
Cohesion: 0.50
Nodes (5): inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), TestCountInodeRows()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.26
Nodes (5): TestFilterExistingWhitelistSymlinkEscape(), applyNewFileMeta(), listEntries(), writeFileInPlace(), vaultFS

### Community 201 - "newTestSession"
Cohesion: 0.50
Nodes (3): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity()

## Knowledge Gaps
- **147 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+142 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 485 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `integrationtests/networkguard_test.go`, `IntegrationSuite`, `IntegrationSuite`, `.runAttrOp`, `.startGuardStd`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `GuardInodeKey`, `BinaryEntry`, `time.Time`, `time.Duration`, `BinaryRule`, `runGuard`, `filevault.go`, `fakeDaemon`, `buildOneGuard`, `sync.Mutex`, `GuardEvent`, `DaemonEvent`, `os.File`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Why does `fileEditModel` connect `fileEditModel` to `catalogrefresh.go`, `fileedit_model_test.go`, `vaultFS`, `Finder`, `fileedit_symlink_test.go`, `highlight_test.go`, `github.com/charmbracelet/bubbletea.Cmd`, `github.com/charmbracelet/bubbles/textarea.Model`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _147 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1185797171198631 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02654867256637168 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.05165289256198347 - nodes in this community are weakly interconnected._