# Graph Report - app-listener  (2026-10-14)

## Corpus Check
- 372 files · ~1,566,107 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4686 nodes · 19508 edges · 188 communities (164 shown, 24 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1644 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `30bd0a8e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- GenerateConf
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- IntegrationSuite
- LibraryClosure
- NetworkGuardUseCase
- go_pkg_path_filepath
- install.go
- updateCatalogConfig
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- Hash
- runMonitor
- guard.bpf.c
- NetworkMonitor
- expandPlaceholders
- inodeChain
- .FilterExistingWhitelist
- uninstall_system.go
- serve.go
- bpfstats/main.go
- catalogRefresher
- fileedit_model_test.go
- guardModel
- time.Duration
- Candidate
- controlServer
- BinaryEntry
- find.go
- main_test.go
- absPath
- Monitor
- Run
- update_test.go
- codeedit.go
- stdio
- golang.org/x/sys/unix.Stat_t
- IntegrationSuite
- mayUpdate
- inode_dev
- netModel
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Model
- TrustGuard
- hashCache
- wipe
- GuardInodeKey
- unit.go
- buildOneGuard
- formatEntry
- exe_supersede.h
- GuardEvent
- runLockdown
- pinstate.go
- is_event_type_allowed
- CandidateDir
- changelog.go
- filevault.go
- fcntl
- copytree.go
- fileedit_symlink_test.go
- time.Time
- btrfsLayoutFrom
- loadDaemonConfig
- github.com/spf13/cobra.Command
- .refreshOnce
- process_vm_readv.c
- bottomBar
- IntegrationSuite
- netGuardModel
- backups.go
- buildTrustedSet
- Ledger
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- mustSubkey
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- IntegrationSuite
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- CleanupStalePins
- KernelDev
- IntegrationSuite
- Well-known bypass classes
- openBunTmp
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- GlobReservations
- Build & sign release assets (reusable workflow)
- eventUnitTest
- runDaemon
- lockAndDeprovision
- launch_rule_at
- app-listener Terminal Demo Recording
- IntegrationSuite
- highlight.go
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- Open
- NewDumpHook
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- writeWithin
- make test-integration (rootful Docker suite)
- bufio.Reader
- supersede_probe.c
- liveSession
- trace-app-libs.sh
- IntegrationSuite
- newWatchPattern
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- EventType
- BinaryStat
- Config
- IntegrationSuite
- Guard
- net_tester/main.go
- NewEventFanout
- newTestVetter
- Vault
- io.Writer
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- lockRootRecovering
- render.go
- highlight_test.go
- tempGrantSetup
- Installation
- os.File
- graphify-refresh.sh
- RawTarget
- io.Reader
- runLiveEdit
- launchKey
- find_test.go
- markReady
- Vault
- ParseEventsFlag
- Limitations
- Troubleshooting
- Daemon configuration
- edit-protected
- Compatibility
- How it works
- bpfLsmListed
- refreshAdmits
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- runNetworkGuard
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS

## God Nodes (most connected - your core abstractions)
1. `Guard` - 114 edges
2. `Config` - 95 edges
3. `GuardInodeKey` - 94 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 66 edges
9. `fileEditModel` - 60 edges
10. `engine` - 58 edges

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

## Communities (188 total, 24 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (29): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+21 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (126): TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffMergePreservesExistingAndParses(), TestDiffCatalogFlagRegistered(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf() (+118 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "GenerateConf"
Cohesion: 0.06
Nodes (52): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, Section, SectionScan, findSectionEnd(), findSectionStart() (+44 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (77): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+69 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.06
Nodes (77): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), promptEditPassword() (+69 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (17): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+9 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.26
Nodes (24): buildGlobReservations(), bitOf(), discordLibConfig(), hasChild(), mkdirs(), steamLibDirConfig(), TestBuildGlobReservations_BunRootsListed(), TestBuildGlobReservations_BunTmpdir() (+16 more)

### Community 10 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 11 - "LibraryClosure"
Cohesion: 0.12
Nodes (17): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+9 more)

### Community 12 - "NetworkGuardUseCase"
Cohesion: 0.15
Nodes (12): NetworkGuardRepository, NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+4 more)

### Community 13 - "go_pkg_path_filepath"
Cohesion: 0.33
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 14 - "install.go"
Cohesion: 0.07
Nodes (55): applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), encryptDirectories(), ensureInstalledBinary(), ensureMasterKey(), installBinaryOnly(), parseMaintenanceFlags() (+47 more)

### Community 15 - "updateCatalogConfig"
Cohesion: 0.15
Nodes (17): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+9 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.18
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 19 - "fileEditModel"
Cohesion: 0.11
Nodes (8): clipLabel(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.16
Nodes (45): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+37 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (20): ownerPaths(), libBinaries(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(), TestMayUpdate_HardLinkWaivedForUserWritableLibBinary() (+12 more)

### Community 23 - "Hash"
Cohesion: 0.13
Nodes (28): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+20 more)

### Community 24 - "runMonitor"
Cohesion: 0.11
Nodes (12): runHeadless(), runMonitor(), runTUI(), FileEvent, Run(), MonitorRepository, NewModel(), MonitorUseCase (+4 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (59): add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial() (+51 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.07
Nodes (18): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), eventsetSummary() (+10 more)

### Community 27 - "expandPlaceholders"
Cohesion: 0.13
Nodes (9): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), TestPathsFor(), LibDirGlob, TrustGlob, reservedGlobs() (+1 more)

### Community 28 - "inodeChain"
Cohesion: 0.21
Nodes (16): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+8 more)

### Community 29 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 30 - "uninstall_system.go"
Cohesion: 0.05
Nodes (63): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), confirmRunPrereq(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep() (+55 more)

### Community 31 - "serve.go"
Cohesion: 0.09
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - "catalogRefresher"
Cohesion: 0.16
Nodes (7): drain(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.15
Nodes (25): fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod() (+17 more)

### Community 35 - "guardModel"
Cohesion: 0.08
Nodes (20): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), listenForGuardEvents(), tickGuardStats(), clipToWidth() (+12 more)

### Community 36 - "time.Duration"
Cohesion: 0.19
Nodes (10): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately() (+2 more)

### Community 37 - "Candidate"
Cohesion: 0.11
Nodes (28): appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), pathCovered(), runDiffCatalog(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates() (+20 more)

### Community 38 - "controlServer"
Cohesion: 0.15
Nodes (10): readConfigPut(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+2 more)

### Community 39 - "BinaryEntry"
Cohesion: 0.12
Nodes (20): HashStore, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), BinariesSummary(), canonicalPaths(), confinedEntry(), claimHash() (+12 more)

### Community 40 - "find.go"
Cohesion: 0.21
Nodes (8): findButton, span, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 41 - "main_test.go"
Cohesion: 0.05
Nodes (20): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+12 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 43 - "Monitor"
Cohesion: 0.06
Nodes (11): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace(), Monitor (+3 more)

### Community 44 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 45 - "update_test.go"
Cohesion: 0.05
Nodes (36): assetsFor(), checkUpdatePreconditions(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), fetchReleases(), parseChecksum(), parsePublicKey() (+28 more)

### Community 46 - "codeedit.go"
Cohesion: 0.25
Nodes (16): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+8 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 50 - "mayUpdate"
Cohesion: 0.33
Nodes (6): hardLinked(), mayUpdate(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), namedGeneralTool(), ResolvesToGeneralTool()

### Community 51 - "inode_dev"
Cohesion: 0.16
Nodes (15): ctx_ptr(), inode_dev(), sb_dev(), current_is_inspector(), exe_tag_bits(), get_current_exe_inode(), get_task_exe_inode(), guard_bprm_check_security() (+7 more)

### Community 52 - "netModel"
Cohesion: 0.12
Nodes (16): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+8 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.15
Nodes (36): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+28 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.12
Nodes (7): renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 55 - "TrustGuard"
Cohesion: 0.05
Nodes (20): tempGrant, trustHook, trustRows, deleteResKeys(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust(), liveTrust() (+12 more)

### Community 56 - "hashCache"
Cohesion: 0.18
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 57 - "wipe"
Cohesion: 0.17
Nodes (8): migrationTarget, openFileVaultWithMasterKey(), writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), wipe(), OpenRegularNoFollow()

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (51): engine, refusedReplacements, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), appendInodeRow() (+43 more)

### Community 59 - "unit.go"
Cohesion: 0.09
Nodes (29): TestDaemonUnitWithMetadataOutput(), installFile(), installFileAs(), installUnit(), upsertFile(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeSSHAgentUnit() (+21 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.06
Nodes (55): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), drainEphemeral(), eventFilterOptions(), newPinGeneration(), reloadOnce() (+47 more)

### Community 61 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 62 - "exe_supersede.h"
Cohesion: 0.16
Nodes (19): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+11 more)

### Community 63 - "GuardEvent"
Cohesion: 0.05
Nodes (20): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), commMatchesGuardedBinary(), decodeGuardEvent() (+12 more)

### Community 64 - "runLockdown"
Cohesion: 0.15
Nodes (15): configureDaemonLogging(), confirmMasterKeyOverwrite(), runGenKey(), runLockdown(), runOneShotMode(), writeInPlace(), readTempJournal(), stripTempJournal() (+7 more)

### Community 65 - "pinstate.go"
Cohesion: 0.29
Nodes (14): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+6 more)

### Community 66 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 67 - "CandidateDir"
Cohesion: 0.30
Nodes (4): entryWriters(), underResource(), globBuilder, CandidateDir

### Community 68 - "changelog.go"
Cohesion: 0.10
Nodes (34): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), showChangelog(), rel(), tags() (+26 more)

### Community 69 - "filevault.go"
Cohesion: 0.27
Nodes (13): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), recoverFileInPlace(), rejectSidecarSymlink(), sealFileVaultWithMasterKey() (+5 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copytree.go"
Cohesion: 0.20
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "fileedit_symlink_test.go"
Cohesion: 0.35
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 73 - "time.Time"
Cohesion: 0.12
Nodes (16): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+8 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 78 - ".refreshOnce"
Cohesion: 0.21
Nodes (11): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), RefreshCatalog(), TestRefreshCatalogAdmitVetsOnlyNewLines() (+3 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 82 - "netGuardModel"
Cohesion: 0.24
Nodes (5): NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 83 - "backups.go"
Cohesion: 0.13
Nodes (24): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), Delete(), Find(), find(), Backup (+16 more)

### Community 84 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 85 - "Ledger"
Cohesion: 0.15
Nodes (3): Ledger, Pending, Source

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.13
Nodes (12): lineStyles, viewportState, column(), Finder, baseStyle(), gutterLabel(), lineLen(), markStyle() (+4 more)

### Community 88 - "mustSubkey"
Cohesion: 0.27
Nodes (10): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering() (+2 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.07
Nodes (36): UIDResolver, NewUIDResolver(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI(), writeEvent() (+28 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.11
Nodes (8): eventTypeKey(), mcNameKey(), statInodeKey(), parseProcs(), TestParseProcs(), hook, NetGuard, procRec

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 100 - "KernelDev"
Cohesion: 0.09
Nodes (24): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs(), KernelDev(), BtrfsDevices() (+16 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "GlobReservations"
Cohesion: 0.12
Nodes (16): reserveChain(), GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, ParseGlobName(), resolveBits() (+8 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "eventUnitTest"
Cohesion: 0.13
Nodes (6): harnessSuite, BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent(), editorHarness

### Community 110 - "runDaemon"
Cohesion: 0.12
Nodes (13): newControlManager(), startControlManager(), awaitStartupOrSignal(), catchLifecycleSignals(), runDaemon(), startGuardedDaemonAbortable(), startPprof(), writePidFile() (+5 more)

### Community 111 - "lockAndDeprovision"
Cohesion: 0.21
Nodes (8): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision(), modifiedContextWithSource()

### Community 112 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "highlight.go"
Cohesion: 0.11
Nodes (8): diffModel, editorModel, fg(), Highlighter, lexerFor(), shebangLexer(), TestTokenizeKeepsLineCount(), tokenClassOf()

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "Open"
Cohesion: 0.20
Nodes (12): ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders(), ensureHashFilePlaceholder(), EnsurePlaceholders(), JournalPath(), Open(), TestEnsurePlaceholdersNoInstallDir() (+4 more)

### Community 121 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

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

### Community 127 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "bufio.Reader"
Cohesion: 0.22
Nodes (10): TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest(), TestReadGrantRequestForward(), TestReadGrantRequestSelect() (+2 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "liveSession"
Cohesion: 0.18
Nodes (12): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries() (+4 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "newWatchPattern"
Cohesion: 0.31
Nodes (11): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+3 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "EventType"
Cohesion: 0.09
Nodes (22): EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), TempBinary, TemporaryGrant, TempRule, planTempBlock() (+14 more)

### Community 143 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 144 - "Config"
Cohesion: 0.09
Nodes (33): newCatalogRefresher(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard(), trustStartupError() (+25 more)

### Community 146 - "Guard"
Cohesion: 0.05
Nodes (16): binaryVerifyState, deferredBinary, exeRow, rootHandle, tempMaskOp, tempRow, openRoot(), SharedPinDegraded() (+8 more)

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 149 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 150 - "Vault"
Cohesion: 0.19
Nodes (10): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+2 more)

### Community 151 - "io.Writer"
Cohesion: 0.36
Nodes (7): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), candidate

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (18): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+10 more)

### Community 154 - "lockRootRecovering"
Cohesion: 0.29
Nodes (8): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), PinPrefix(), sanitizeGen(), SharedPinPrefix(), TestPinPrefixDistinctAndStable(), TestPinPrefixShape()

### Community 156 - "highlight_test.go"
Cohesion: 0.20
Nodes (22): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+14 more)

### Community 157 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.06
Nodes (51): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), setConfinedHomes(), binaryVetter, vetDecision, VettedInode, Multicall (+43 more)

### Community 161 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 162 - "io.Reader"
Cohesion: 0.33
Nodes (5): BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype()

### Community 163 - "runLiveEdit"
Cohesion: 0.40
Nodes (5): pingInterval(), chooseResource(), runLiveEdit(), RunFileEditorSession(), RunSession()

### Community 165 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 167 - "markReady"
Cohesion: 0.25
Nodes (9): markReady(), notifySystemdReady(), inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), InodeOccupancy(), TestCountInodeRows() (+1 more)

### Community 171 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 185 - "refreshAdmits"
Cohesion: 0.50
Nodes (4): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), BinaryRule

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "runNetworkGuard"
Cohesion: 0.12
Nodes (23): CheckBPFLSM(), CheckEBPF(), ParseNetEventsFlag(), prepareDaemonStart(), runBPFCheck(), computeGuardBinaries(), confirmUnsafe(), discoverInfra() (+15 more)

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

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.11
Nodes (17): checkFreshDir(), openDirNoFollow(), renewDir(), AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode(), humanSize(), isBinaryContent() (+9 more)

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 451 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **24 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `time.Duration`, `BinaryEntry`, `guardUnitTest`, `time.Time`, `runDaemon`, `DaemonEvent`, `TrustGuard`, `GuardInodeKey`, `buildOneGuard`, `GuardEvent`, `os.File`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11271848946835311 - nodes in this community are weakly interconnected._
- **Why does `engine` connect `GuardInodeKey` to `bpfstats/main.go`, `daemon/daemon.go`, `Guard`, `TrustGuard`, `runNetworkGuard`, `NetGuard`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.025221456692913386 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08982456140350877 - nodes in this community are weakly interconnected._