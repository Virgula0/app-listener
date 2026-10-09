# Graph Report - app-listener  (2026-10-09)

## Corpus Check
- 329 files · ~1,515,340 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4403 nodes · 17743 edges · 172 communities (143 shown, 29 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1290 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4f253dde`
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
- mkdirs
- IntegrationSuite
- migrate.go
- codeedit.go
- runUpdate
- install.go
- Vault
- IntegrationSuite
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- runMonitor
- os.File
- guard.bpf.c
- app.go
- github.com/spf13/cobra.Command
- reloadOnce
- buildOneGuard
- uninstall_system.go
- Serve
- engine
- IntegrationSuite
- fileedit_model_test.go
- DaemonEvent
- watchGroup
- New
- GuardEvent
- guardSpec
- github.com/charmbracelet/bubbles/textarea.Model
- NetworkMonitor
- absPath
- time.Duration
- IntegrationSuite
- fscrypt/orphans.go
- Monitor
- .syncSetsLocked
- IntegrationSuite
- expandPlaceholders
- __always_inline
- IntegrationSuite
- MountVouch
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- GuardInodeKey
- controlServer
- fscrypt.go
- EventType
- TrustGuard
- runNetworkGuard
- exe_supersede.h
- SanitizeText
- errno
- pinstate.go
- User
- highlight.go
- LibraryClosure
- catalogRefresher
- fcntl
- copyRegularAt
- update_test.go
- buildTrustedSet
- btrfsLayoutFrom
- Ledger
- IntegrationSuite
- mayUpdate
- process_vm_readv.c
- guardModel
- uninstall.go
- GlobReservations
- NewEventFanout
- .step
- system.go
- Highlighter
- time.Time
- TempBinary
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- KernelDev
- CandidateDir
- netinject.c
- bufio.Reader
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- RawTarget
- inodeChain
- IntegrationSuite
- Well-known bypass classes
- Vault
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- NetGuard
- Build & sign release assets (reusable workflow)
- globBuilder
- Candidate
- net_tester/main.go
- TestMonitorUseCaseLifecycle
- app-listener Terminal Demo Recording
- NetworkMonitorUseCase
- is_event_type_allowed
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- testKeyPair
- bottomBar
- WebSocket /ws client
- NewDumpHook
- launch_probe.c
- install.sh
- mc_tag.h
- Config
- make test-integration (rootful Docker suite)
- vaultFS
- supersede_probe.c
- launchKey
- trace-app-libs.sh
- IntegrationSuite
- controlManager
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- netGuardModel
- find_test.go
- IntegrationSuite
- Guard
- tempRow
- loadDaemonConfig
- configedit_test.go
- BinaryEntry
- launch_rule_at
- bpfLsmListed
- update.go
- downloadAndVerify
- tempGrantSetup
- planTaintOwners
- TestMain
- sync.Mutex
- .runEditorHarness
- graphify-refresh.sh
- binaryVerifyState
- runLockdown
- newChangelogModel
- Vault
- fakeDaemon
- io.Reader
- Discover
- ParseSectionWhitelist
- .SystemWhitelist
- TestCatalogTrustGlobsAreReservable

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

## Communities (172 total, 29 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (31): atomicWriteAt(), writeAndSync(), Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent (+23 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (119): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+111 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "filevault.go"
Cohesion: 0.14
Nodes (28): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, looksLikeFileVaultRecord(), newFileVaultAEAD() (+20 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (74): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+66 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.18
Nodes (31): rcTarget, bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), ensureRCBlock(), EnsureSSHAgentEnv(), openOwnedRegular(), rcBlock() (+23 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (20): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+12 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "IntegrationSuite"
Cohesion: 0.06
Nodes (46): ParseEventsFlag(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+38 more)

### Community 12 - "codeedit.go"
Cohesion: 0.12
Nodes (39): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), moveTo(), Navigate() (+31 more)

### Community 13 - "runUpdate"
Cohesion: 0.15
Nodes (19): changelogText(), showChangelog(), applyUpdate(), confirmUpdate(), filterChannel(), isTerminal(), latestTime(), newerThanInstalled() (+11 more)

### Community 14 - "install.go"
Cohesion: 0.06
Nodes (55): reloadDaemonForPasswordChange(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), encryptDirectories(), ensureInstalledBinary(), installBinaryOnly() (+47 more)

### Community 15 - "Vault"
Cohesion: 0.12
Nodes (18): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+10 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.12
Nodes (7): eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "GenerateConf"
Cohesion: 0.06
Nodes (50): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+42 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (18): backingDeviceUnion(), buildConcurrency(), buildGuards(), libBinaries(), planUpdaters(), resourceBinaries(), resourceLinks(), TestPlanUpdaters_LinkOwnsItsResourceBit() (+10 more)

### Community 23 - "runMonitor"
Cohesion: 0.11
Nodes (15): CheckEBPF(), ParseNetEventsFlag(), prepareDaemonStart(), runHeadless(), runMonitor(), runTUI(), runNetworkMonitor(), runTUI() (+7 more)

### Community 24 - "os.File"
Cohesion: 0.08
Nodes (36): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+28 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (46): add_inode_to_guard(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), evict_inode_from_guard(), fill_path(), get_dentry_from_path(), get_inode_from_path() (+38 more)

### Community 26 - "app.go"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 27 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 28 - "reloadOnce"
Cohesion: 0.11
Nodes (23): lockOneRoot(), lockRootRecovering(), makeReloadHandler(), newPinGeneration(), reloadOnce(), reloadSlotsNeeded(), relockStaleVaults(), startCatalogRefresh() (+15 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.13
Nodes (23): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), selfGuards, selfProtectSpec, AdmissionCheck (+15 more)

### Community 30 - "uninstall_system.go"
Cohesion: 0.29
Nodes (11): detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertSSHAgents(), staleRCUsers() (+3 more)

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (28): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+20 more)

### Community 32 - "engine"
Cohesion: 0.09
Nodes (6): TestTieRank(), engine, pathWithin(), tieRank(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent()

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.07
Nodes (47): TestCleanupStalePins(), TestFilterExistingWhitelistSymlinkEscape(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf() (+39 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.09
Nodes (31): UIDResolver, NewUIDResolver(), drainEphemeral(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter(), targetProgram() (+23 more)

### Community 36 - "watchGroup"
Cohesion: 0.11
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 37 - "New"
Cohesion: 0.17
Nodes (17): harnessSuite, TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip(), TestVaultUnlockLockDispatchForRegularFile() (+9 more)

### Community 38 - "GuardEvent"
Cohesion: 0.07
Nodes (22): displayEntries(), holdForward(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI() (+14 more)

### Community 39 - "guardSpec"
Cohesion: 0.15
Nodes (18): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), InertProgram(), KernelHasFunc(), TestKernelHasFunc(), trustSpec() (+10 more)

### Community 40 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.18
Nodes (9): findButton, span, column(), findAll(), Finder, lowerRunes(), overlay(), runesEqual() (+1 more)

### Community 41 - "NetworkMonitor"
Cohesion: 0.10
Nodes (14): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), eventsetSummary() (+6 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "time.Duration"
Cohesion: 0.21
Nodes (5): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), parseErr(), liveSession

### Community 44 - "IntegrationSuite"
Cohesion: 0.17
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "fscrypt/orphans.go"
Cohesion: 0.23
Nodes (19): appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies() (+11 more)

### Community 46 - "Monitor"
Cohesion: 0.06
Nodes (8): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), Monitor, monitorUnitTest, pathCache, pathCacheEntry

### Community 48 - ".syncSetsLocked"
Cohesion: 0.16
Nodes (11): resInfo, engine, intersectAllows(), isSubset(), planMemberRows(), setPaths(), syncSetMap(), syncU32Map() (+3 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "expandPlaceholders"
Cohesion: 0.25
Nodes (7): symlinkTargetGlob(), expandPlaceholders(), CandidateDir, LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 51 - "__always_inline"
Cohesion: 0.11
Nodes (27): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), btrfs_copy_gate(), check_and_emit_args(), code_suspect() (+19 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 53 - "MountVouch"
Cohesion: 0.13
Nodes (9): deleteInoKeys(), deleteResKeys(), syncMap(), deleteKeys(), MountVouch, TrustGuard, NewMountVouch(), NewTrustGuard() (+1 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (14): diffModel, editorModel, newFindInput(), clamp(), NewModel(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), assertViewFits() (+6 more)

### Community 56 - "GuardInodeKey"
Cohesion: 0.04
Nodes (61): refusedReplacements, engine, Inspector, TrustGuard, SetInspectors(), engine, Guard, supersedeUnnamed() (+53 more)

### Community 57 - "controlServer"
Cohesion: 0.18
Nodes (8): controlManager, controlServer, readClientLines(), startControlServer(), streamEvents(), configEditor, editControlSession, grantRequest

### Community 58 - "fscrypt.go"
Cohesion: 0.10
Nodes (24): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), ensureMasterKey(), checkKeyLen(), classifySetupError(), classifySupportError(), GenerateMasterKey() (+16 more)

### Community 59 - "EventType"
Cohesion: 0.09
Nodes (14): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock() (+6 more)

### Community 60 - "TrustGuard"
Cohesion: 0.13
Nodes (6): trustHook, trustRows, setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 61 - "runNetworkGuard"
Cohesion: 0.13
Nodes (16): CheckBPFLSM(), runBPFCheck(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI() (+8 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (17): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+9 more)

### Community 63 - "SanitizeText"
Cohesion: 0.06
Nodes (16): multicallAdmissible(), ownerPaths(), systemPatterns(), systemRule(), runHeadless(), runHeadless(), binaryVetter, vetDecision (+8 more)

### Community 64 - "errno"
Cohesion: 0.21
Nodes (5): put(), wait_then_read(), dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (23): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+15 more)

### Community 66 - "User"
Cohesion: 0.06
Nodes (53): trustManager, logRefreshChange(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), TestLogRefreshChange(), askBunTmpdirUsers(), bunEntryConfigured() (+45 more)

### Community 67 - "highlight.go"
Cohesion: 0.23
Nodes (5): fg(), Highlighter, lexerFor(), shebangLexer(), tokenClassOf()

### Community 68 - "LibraryClosure"
Cohesion: 0.10
Nodes (19): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+11 more)

### Community 69 - "catalogRefresher"
Cohesion: 0.19
Nodes (5): drain(), newCatalogRefresher(), readEvents(), catalogRefresher, debouncer

### Community 70 - "fcntl"
Cohesion: 0.13
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "update_test.go"
Cohesion: 0.15
Nodes (3): fetchReleases(), TestFetchReleases(), TestFetchReleasesHTTPError()

### Community 73 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "Ledger"
Cohesion: 0.09
Nodes (25): newBinaryVetter(), newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved() (+17 more)

### Community 78 - "mayUpdate"
Cohesion: 0.20
Nodes (11): hardLinked(), mayUpdate(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), TestMayUpdate_HardLinkWaivedForUserWritableLibBinary(), TestPlanUpdaters_ScopedToOwnResources() (+3 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "guardModel"
Cohesion: 0.05
Nodes (36): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+28 more)

### Community 81 - "uninstall.go"
Cohesion: 0.07
Nodes (44): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), deletePostBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories() (+36 more)

### Community 82 - "GlobReservations"
Cohesion: 0.17
Nodes (10): GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, TrustGuard, resolveBits(), resolveChildren() (+2 more)

### Community 83 - "NewEventFanout"
Cohesion: 0.24
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 84 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 85 - "system.go"
Cohesion: 0.08
Nodes (37): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+29 more)

### Community 86 - "Highlighter"
Cohesion: 0.15
Nodes (13): lineStyles, viewportState, baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 87 - "time.Time"
Cohesion: 0.14
Nodes (14): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+6 more)

### Community 88 - "TempBinary"
Cohesion: 0.10
Nodes (18): TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows(), applyTempGrants(), closeTempBinaries() (+10 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "KernelDev"
Cohesion: 0.07
Nodes (35): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding() (+27 more)

### Community 92 - "CandidateDir"
Cohesion: 0.18
Nodes (11): confSafeMatch(), BinaryRule, CandidateDir, homeMatchConfined(), symlinkStaysInParent(), boundedGlob(), exists(), matchDir() (+3 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 94 - "bufio.Reader"
Cohesion: 0.14
Nodes (13): swapFile(), controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines() (+5 more)

### Community 95 - "NetGuard"
Cohesion: 0.12
Nodes (8): FirstOnBtrfs(), OnBtrfs(), SuperblockDevPath(), TestSuperblockDevMatchesStatOffBtrfs(), eventTypeKey(), NetGuard, mcNameKey(), statInodeKey()

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 100 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "Vault"
Cohesion: 0.26
Nodes (6): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata()

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

### Community 109 - "globBuilder"
Cohesion: 0.30
Nodes (6): entryWriters(), reserveChain(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName()

### Community 110 - "Candidate"
Cohesion: 0.26
Nodes (14): appendSectionsAndEdit(), collectDiffAdditions(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories() (+6 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (11): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 121 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "Config"
Cohesion: 0.06
Nodes (46): openBinaryVetter(), sameConfig(), awaitStartupOrSignal(), catchLifecycleSignals(), newDaemonModel(), notifySystemdReady(), runDaemon(), runDaemonUI() (+38 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "vaultFS"
Cohesion: 0.31
Nodes (4): applyNewFileMeta(), listEntries(), dirEntry, vaultFS

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "controlManager"
Cohesion: 0.29
Nodes (3): controlManager, newControlManager(), startControlManager()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "netGuardModel"
Cohesion: 0.26
Nodes (5): NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 143 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 146 - "Guard"
Cohesion: 0.06
Nodes (18): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, rootHandle, openRoot(), SharedPinDegraded() (+10 more)

### Community 147 - "tempRow"
Cohesion: 0.18
Nodes (6): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust()

### Community 148 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), SetUserHomes(), TestResolveConfinedRootOwnedParents(), TestOpenSystemPlacedRefusesUserHome()

### Community 149 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 150 - "BinaryEntry"
Cohesion: 0.18
Nodes (10): BinaryEntry, checkBtrfsKeys(), BinariesSummary(), canonicalPaths(), commMatchesGuardedBinary(), BinariesSummary(), entryOf(), ComputeBinaryEntryFile() (+2 more)

### Community 151 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 152 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 153 - "update.go"
Cohesion: 0.31
Nodes (7): compareStableVersions(), newerThanStable(), parseStableVersion(), parseVerPart(), TestNewerThanStable(), TestParseStableVersion(), stableVersion

### Community 154 - "downloadAndVerify"
Cohesion: 0.31
Nodes (7): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), sanityCheckBinary(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), updateHTTPClient()

### Community 155 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 156 - "planTaintOwners"
Cohesion: 0.31
Nodes (8): taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), planTaintOwners(), setKeyOf(), allowSet(), TestPlanTaintOwners(), TestPlanTaintOwners_ReloadOverlapKeepsSetKey()

### Community 157 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 161 - "binaryVerifyState"
Cohesion: 0.43
Nodes (3): binaryVerifyState, sharedHashEntry, hashBinaryShared()

### Community 162 - "runLockdown"
Cohesion: 0.29
Nodes (8): configureDaemonLogging(), resolveConfigPath(), runLockdown(), dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 163 - "newChangelogModel"
Cohesion: 0.43
Nodes (6): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes()

### Community 166 - "io.Reader"
Cohesion: 0.33
Nodes (4): BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype(), progressReader

### Community 167 - "Discover"
Cohesion: 0.53
Nodes (4): Discover(), DiscoverSystem(), TestDiscoverOnlyExisting(), TestDiscoverSystemOnlyAbsolute()

### Community 168 - "ParseSectionWhitelist"
Cohesion: 0.33
Nodes (6): isLibDirectiveText(), ParseSectionHeaderPath(), unquotePath(), IsLibraryDirective(), ParseSectionWhitelist(), TestParseSectionWhitelist()

## Knowledge Gaps
- **53 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 386 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **29 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`?**
  _High betweenness centrality (0.071) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `binaryVerifyState`, `DaemonEvent`, `fakeDaemon`, `controlManager`, `GuardEvent`, `BinaryEntry`, `time.Time`, `os.File`, `GuardInodeKey`, `reloadOnce`, `buildOneGuard`, `sync.Mutex`, `SanitizeText`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **Why does `Vault` connect `Vault` to `New`, `install.go`, `uninstall.go`, `GenerateConf`, `Guard`, `system.go`, `fscrypt.go`, `reloadOnce`, `sync.Mutex`, `Config`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _53 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10476943346508565 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02438693943910039 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._