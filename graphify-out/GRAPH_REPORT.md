# Graph Report - app-listener  (2026-10-09)

## Corpus Check
- 329 files · ~1,517,164 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4409 nodes · 17768 edges · 166 communities (143 shown, 23 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1294 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0def5117`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- filevault.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- auth.go
- migrate.go
- codeedit.go
- runUpdate
- systemd.go
- Vault
- IntegrationSuite
- guard_trust.bpf.c
- confgen.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- runMonitor
- openBunTmp
- guard.bpf.c
- Run
- github.com/spf13/cobra.Command
- runDaemon
- buildOneGuard
- CandidateDir
- serve.go
- engine
- IntegrationSuite
- fileedit_model_test.go
- DaemonEvent
- watchGroup
- New
- GuardEvent
- bpfstats/main.go
- Finder
- NetworkMonitor
- IntegrationSuite
- liveSession
- IntegrationSuite
- fscrypt/orphans.go
- Monitor
- .syncSetsLocked
- IntegrationSuite
- expandPlaceholders
- __always_inline
- IntegrationSuite
- TrustGuard
- github.com/charmbracelet/bubbletea.Model
- absPath
- os.File
- controlServer
- install.go
- EventType
- fscrypt.go
- runNetworkGuard
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- IntegrationSuite
- highlight_test.go
- LibraryClosure
- catalogRefresher
- fcntl
- copytree.go
- bufio.Reader
- GlobReservations
- btrfsLayoutFrom
- Ledger
- inode_dev
- planUpdaters
- process_vm_readv.c
- netModel
- backups.go
- update_test.go
- globBuilder
- .step
- system.go
- github.com/charmbracelet/bubbles/textarea.Model
- watchSet
- TempBinary
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- KernelDev
- runForward
- netinject.c
- time.Duration
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- uninstall_system.go
- newWatchPattern
- IntegrationSuite
- Well-known bypass classes
- Vault
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- NetGuard
- Build & sign release assets (reusable workflow)
- newTestVetter
- Candidate
- net_tester/main.go
- TestMonitorUseCaseLifecycle
- app-listener Terminal Demo Recording
- is_event_type_allowed
- update.go
- lockAndDeprovision
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- mustSubkey
- testcontainers.Container
- WebSocket /ws client
- NewDumpHook
- launch_probe.c
- install.sh
- mc_tag.h
- Config
- make test-integration (rootful Docker suite)
- time.Time
- supersede_probe.c
- lockRootRecovering
- trace-app-libs.sh
- IntegrationSuite
- sync.Mutex
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- bottomBar
- find_test.go
- io.Writer
- Guard
- GuardInodeKey
- downloadAndVerify
- configedit_test.go
- ParseEventsFlag
- writeWithin
- AddServeFlags
- runEditConfig
- formatEntry
- tempGrantSetup
- launchKey
- graphify-refresh.sh
- downloadFile
- prepareDaemonStart
- changelog.go
- DiscoverForUsers
- launch_rule_at
- Vault
- .SystemWhitelist

## God Nodes (most connected - your core abstractions)
1. `Config` - 93 edges
2. `Guard` - 93 edges
3. `GuardInodeKey` - 86 edges
4. `Load()` - 85 edges
5. `absPath()` - 78 edges
6. `IntegrationSuite` - 75 edges
7. `IntegrationSuite` - 66 edges
8. `writeConfig()` - 65 edges
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

## Communities (166 total, 23 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (36): atomicWriteAt(), writeAndSync(), Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent (+28 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (117): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+109 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.06
Nodes (13): IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe() (+5 more)

### Community 3 - "filevault.go"
Cohesion: 0.21
Nodes (18): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, looksLikeFileVaultRecord(), openFileVaultWithMasterKey() (+10 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+70 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.10
Nodes (50): trustManager, askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), deploy(), preflightDeployedBinary() (+42 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (26): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf() (+18 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.26
Nodes (24): buildGlobReservations(), bitOf(), discordLibConfig(), hasChild(), mkdirs(), steamLibDirConfig(), TestBuildGlobReservations_BunRootsListed(), TestBuildGlobReservations_BunTmpdir() (+16 more)

### Community 10 - "auth.go"
Cohesion: 0.12
Nodes (27): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+19 more)

### Community 11 - "migrate.go"
Cohesion: 0.34
Nodes (4): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), hint

### Community 12 - "codeedit.go"
Cohesion: 0.18
Nodes (26): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+18 more)

### Community 13 - "runUpdate"
Cohesion: 0.18
Nodes (16): changelogText(), applyUpdate(), confirmUpdate(), fetchReleases(), filterChannel(), isTerminal(), latestTime(), newerThanInstalled() (+8 more)

### Community 14 - "systemd.go"
Cohesion: 0.06
Nodes (54): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), reloadDaemonForPasswordChange(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), TestCollectFilesystemPrereqsNoPanic() (+46 more)

### Community 15 - "Vault"
Cohesion: 0.20
Nodes (11): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+3 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.15
Nodes (5): eventsForPath(), IntegrationSuite, guardDeltaEvents(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "confgen.go"
Cohesion: 0.06
Nodes (53): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, Section, SectionScan, ConfSafePath(), findSectionEnd() (+45 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (15): applyNewFileMeta(), clipLabel(), humanSize(), listEntries(), clamp(), sortNodes(), listenForGuardEvents(), tickGuardStats() (+7 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.10
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "runMonitor"
Cohesion: 0.11
Nodes (19): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+11 more)

### Community 24 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 25 - "guard.bpf.c"
Cohesion: 0.14
Nodes (37): add_inode_to_guard(), chmod_only_drops_write(), discover_guarded_parent(), evict_inode_from_guard(), fill_path(), get_dentry_from_path(), get_inode_from_path(), guard_file_truncate() (+29 more)

### Community 26 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 27 - "github.com/spf13/cobra.Command"
Cohesion: 0.23
Nodes (13): credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+5 more)

### Community 28 - "runDaemon"
Cohesion: 0.14
Nodes (20): backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), makeReloadHandler(), notifySystemdReady(), reloadOnce(), reloadSlotsNeeded() (+12 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.09
Nodes (30): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan() (+22 more)

### Community 30 - "CandidateDir"
Cohesion: 0.23
Nodes (8): confSafeMatch(), BinaryRule, CandidateDir, homeMatchConfined(), symlinkStaysInParent(), boundedGlob(), exists(), matchDir()

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "engine"
Cohesion: 0.09
Nodes (7): TestTieRank(), engine, pathWithin(), tieRank(), deleteInoKeys(), intersectAllows(), TestIntersectAllows()

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.06
Nodes (53): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape() (+45 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.09
Nodes (26): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+18 more)

### Community 36 - "watchGroup"
Cohesion: 0.11
Nodes (33): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+25 more)

### Community 37 - "New"
Cohesion: 0.11
Nodes (27): TestDiffMergePreservesExistingAndParses(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips(), TestGroupedSectionAddressedByEncryptionRoot() (+19 more)

### Community 38 - "GuardEvent"
Cohesion: 0.05
Nodes (29): selfProtectSpecs(), displayEntries(), holdForward(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless() (+21 more)

### Community 39 - "bpfstats/main.go"
Cohesion: 0.17
Nodes (18): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), InertProgram(), KernelHasFunc(), TestKernelHasFunc(), trustSpec() (+10 more)

### Community 40 - "Finder"
Cohesion: 0.14
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 41 - "NetworkMonitor"
Cohesion: 0.07
Nodes (16): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), eventTypeKey() (+8 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.21
Nodes (4): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, netGuardTail()

### Community 45 - "fscrypt/orphans.go"
Cohesion: 0.18
Nodes (21): cleanOrphanedFscrypt(), modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors() (+13 more)

### Community 46 - "Monitor"
Cohesion: 0.06
Nodes (8): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), Monitor, monitorUnitTest, pathCache, pathCacheEntry

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "expandPlaceholders"
Cohesion: 0.25
Nodes (7): symlinkTargetGlob(), expandPlaceholders(), CandidateDir, LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 51 - "__always_inline"
Cohesion: 0.15
Nodes (24): exe_is_superseded(), exe_refused(), check_and_emit_args(), code_suspect(), count_degrade(), current_is_inspector(), emit_process_denial(), emit_with() (+16 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.05
Nodes (17): harnessSuite, IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex() (+9 more)

### Community 53 - "TrustGuard"
Cohesion: 0.06
Nodes (17): trustHook, trustRows, deleteResKeys(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust(), syncU32Map(), cStr() (+9 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.11
Nodes (8): renderBullets(), TestRenderBulletsVerbatimAndIndented(), RunFileEditorSession(), RunSession(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 55 - "absPath"
Cohesion: 0.09
Nodes (8): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath(), TestMain(), helperChild(), TestMain()

### Community 56 - "os.File"
Cohesion: 0.05
Nodes (52): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes() (+44 more)

### Community 57 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 58 - "install.go"
Cohesion: 0.10
Nodes (31): applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), promptEditPassword(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups() (+23 more)

### Community 59 - "EventType"
Cohesion: 0.09
Nodes (15): systemPatterns(), systemRule(), BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+7 more)

### Community 60 - "fscrypt.go"
Cohesion: 0.13
Nodes (19): checkKeyLen(), classifySetupError(), classifySupportError(), GenerateMasterKey(), isNotEncryptedErr(), readKeyFrom(), readOrCreateKey(), TestClassifySetupErrorGeneric() (+11 more)

### Community 61 - "runNetworkGuard"
Cohesion: 0.09
Nodes (20): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), eventsetSummary(), Mode (+12 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 63 - "fakeGuardRepo"
Cohesion: 0.07
Nodes (9): dropSupersededLocked(), Guard, SupersededBinary, PruneSuperseded(), supersededMark(), supersedeMaps(), supersedePruneLoop(), warnLostStamps() (+1 more)

### Community 64 - "errno"
Cohesion: 0.21
Nodes (5): put(), wait_then_read(), dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (23): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+15 more)

### Community 67 - "highlight_test.go"
Cohesion: 0.10
Nodes (20): diffModel, editorModel, Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo() (+12 more)

### Community 68 - "LibraryClosure"
Cohesion: 0.10
Nodes (19): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+11 more)

### Community 69 - "catalogRefresher"
Cohesion: 0.15
Nodes (9): drain(), newCatalogRefresher(), readEvents(), sameConfig(), awaitStartupOrSignal(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup(), TestAwaitStartupOrSignalCompletes() (+1 more)

### Community 70 - "fcntl"
Cohesion: 0.13
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copytree.go"
Cohesion: 0.20
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 73 - "GlobReservations"
Cohesion: 0.19
Nodes (10): GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, TrustGuard, resolveBits(), resolveChildren() (+2 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "Ledger"
Cohesion: 0.12
Nodes (9): EnsurePlaceholders(), Ledger, Pending, Source, JournalPath(), Open(), TestEnsurePlaceholdersNoInstallDir(), TestLedgerKeepsJournalFile() (+1 more)

### Community 76 - "inode_dev"
Cohesion: 0.16
Nodes (12): ctx_ptr(), inode_dev(), sb_dev(), btrfs_copy_gate(), guard_bprm_committed(), guard_exec_applet(), guard_file_ioctl(), guard_file_ioctl_compat() (+4 more)

### Community 78 - "planUpdaters"
Cohesion: 0.15
Nodes (17): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+9 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "netModel"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 81 - "backups.go"
Cohesion: 0.19
Nodes (17): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), Delete(), Find(), find(), Backup (+9 more)

### Community 82 - "update_test.go"
Cohesion: 0.13
Nodes (8): readInstalledVersion(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestReadInstalledVersion(), TestVerifyRelease()

### Community 83 - "globBuilder"
Cohesion: 0.20
Nodes (8): entryWriters(), reserveChain(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), inOwnTree(), TestCatalogTrustGlobsAreReservable()

### Community 84 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 85 - "system.go"
Cohesion: 0.09
Nodes (36): checkAndApply(), pathCovered(), TestPathCovered(), TestDaemonUnitWithMetadataOutput(), isInsidePath(), TestAskSSHAgentUsersNoGuardedSSH(), updateCatalogConfig(), askSSHAgentUsers() (+28 more)

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.16
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "watchSet"
Cohesion: 0.16
Nodes (12): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern (+4 more)

### Community 88 - "TempBinary"
Cohesion: 0.09
Nodes (23): BinaryStat, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow() (+15 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "KernelDev"
Cohesion: 0.06
Nodes (43): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+35 more)

### Community 92 - "runForward"
Cohesion: 0.16
Nodes (14): auditAfterEdit(), dialLiveSession(), putConfig(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), normalizedEvents() (+6 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 94 - "time.Duration"
Cohesion: 0.21
Nodes (4): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer

### Community 95 - "NetGuard"
Cohesion: 0.12
Nodes (9): FirstOnBtrfs(), OnBtrfs(), SuperblockDevPath(), TestSuperblockDevMatchesStatOffBtrfs(), reasonText(), NetGuard, mcNameKey(), statInodeKey() (+1 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "uninstall_system.go"
Cohesion: 0.11
Nodes (24): askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir() (+16 more)

### Community 100 - "newWatchPattern"
Cohesion: 0.31
Nodes (11): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "Vault"
Cohesion: 0.23
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

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

### Community 109 - "newTestVetter"
Cohesion: 0.42
Nodes (9): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSwapWithinGenerationLosesRights() (+1 more)

### Community 110 - "Candidate"
Cohesion: 0.17
Nodes (19): appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig() (+11 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (11): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 115 - "update.go"
Cohesion: 0.22
Nodes (11): compareStableVersions(), newerThanStable(), parseChecksum(), parsePublicKey(), parseStableVersion(), parseVerPart(), TestNewerThanStable(), TestParseStableVersion() (+3 more)

### Community 116 - "lockAndDeprovision"
Cohesion: 0.27
Nodes (7): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 127 - "Config"
Cohesion: 0.07
Nodes (44): swapFile(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections() (+36 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "time.Time"
Cohesion: 0.18
Nodes (14): newWatchSet(), newGateLogLimiter(), gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters() (+6 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "lockRootRecovering"
Cohesion: 0.29
Nodes (8): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), PinPrefix(), sanitizeGen(), SharedPinPrefix(), TestPinPrefixDistinctAndStable(), TestPinPrefixShape()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 143 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 144 - "io.Writer"
Cohesion: 0.43
Nodes (7): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), candidate

### Community 146 - "Guard"
Cohesion: 0.06
Nodes (26): binaryVerifyState, deferredBinary, rootHandle, sharedHashEntry, BinaryEntry, openRoot(), checkBtrfsKeys(), SharedPinDegraded() (+18 more)

### Community 147 - "GuardInodeKey"
Cohesion: 0.07
Nodes (24): exeRow, refusedReplacements, tempGrant, tempMaskOp, tempRow, engine, Inspector, SetInspectors() (+16 more)

### Community 148 - "downloadAndVerify"
Cohesion: 0.15
Nodes (15): assetsFor(), checkUpdatePreconditions(), downloadAndVerify(), downloadReleaseFiles(), resolveAssets(), sanityCheckBinary(), TestAssetsFor(), TestResolveAssetsPicksArch() (+7 more)

### Community 149 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 150 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 151 - "writeWithin"
Cohesion: 0.29
Nodes (7): descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeWithin()

### Community 152 - "AddServeFlags"
Cohesion: 0.33
Nodes (6): AddServeFlags(), init(), init(), init(), init(), init()

### Community 153 - "runEditConfig"
Cohesion: 0.40
Nodes (6): applyEdit(), editAgain(), fetchConfig(), nextConfig(), runEditConfig(), EditText()

### Community 154 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 155 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 156 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 161 - "downloadFile"
Cohesion: 0.17
Nodes (9): downloadFile(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype() (+1 more)

### Community 162 - "prepareDaemonStart"
Cohesion: 0.08
Nodes (28): CheckBPFLSM(), CheckEBPF(), ParseNetEventsFlag(), configureDaemonLogging(), confirmMasterKeyOverwrite(), loadDaemonConfig(), newPinGeneration(), prepareDaemonStart() (+20 more)

### Community 163 - "changelog.go"
Cohesion: 0.22
Nodes (8): newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 164 - "DiscoverForUsers"
Cohesion: 0.33
Nodes (7): Discover(), DiscoverForUsers(), DiscoverSystem(), TestDiscoverForUsersDeduplicates(), TestDiscoverOnlyExisting(), TestDiscoverPerUserHomes(), TestDiscoverSystemOnlyAbsolute()

### Community 166 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

## Knowledge Gaps
- **53 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 385 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `time.Time`, `DaemonEvent`, `sync.Mutex`, `GuardEvent`, `GuardInodeKey`, `os.File`, `runDaemon`, `buildOneGuard`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Why does `newEventTypes()` connect `IntegrationSuite` to `daemon/daemon.go`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _53 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10280960602731125 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.025210084033613446 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.06158415841584158 - nodes in this community are weakly interconnected._