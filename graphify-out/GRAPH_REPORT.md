# Graph Report - app-listener  (2026-10-16)

## Corpus Check
- 376 files · ~1,569,992 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4741 nodes · 19702 edges · 199 communities (168 shown, 31 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1662 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `56e753f5`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- GuardInodeKey
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- mkdirs
- BinaryEntry
- liveSession
- go_pkg_path_filepath
- Config
- install.go
- uninstall_system.go
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- IntegrationSuite
- guard.bpf.c
- NetworkMonitor
- highlight_test.go
- LibraryClosure
- IntegrationSuite
- orphans_test.go
- serve.go
- bpfstats/main.go
- github.com/stretchr/testify/suite.Suite
- fileedit_model_test.go
- Ledger
- uninstall.go
- Monitor
- controlServer
- ConfinedEntry
- releaseAssets
- daemon_multicall_test.go
- absPath
- monitorUnitTest
- newWatchPattern
- update_test.go
- expandPlaceholders
- stdio
- prepareDaemonStart
- .FilterExistingWhitelist
- planUpdaters
- codeedit.go
- Read
- watchGroup
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- system.go
- IntegrationSuite
- engine
- github.com/charmbracelet/bubbles/textarea.Model
- CandidateDir
- fileedit_symlink_test.go
- exe_supersede.h
- GuardEvent
- netGuardModel
- pinstate.go
- gui/session.go
- runNetworkMonitor
- changelog.go
- buildTrustedSet
- fcntl
- copyRegularAt
- KernelDev
- time.Time
- btrfsLayoutFrom
- confgen.go
- github.com/spf13/cobra.Command
- Highlighter
- process_vm_readv.c
- net_tester/main.go
- IntegrationSuite
- netModel
- time.Duration
- Pulse
- downloadAndVerify
- auditGroupCoverage
- fakeGuardRepo
- unit.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- catalogRefresher
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- changelog_test.go
- golang.org/x/sys/unix.Stat_t
- IntegrationSuite
- Well-known bypass classes
- writeWithin
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- NetGuardEvent
- Build & sign release assets (reusable workflow)
- EventType
- fakeDaemon
- sync.Mutex
- testKeyPair
- app-listener Terminal Demo Recording
- IntegrationSuite
- .PathsFor
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- compileGlobNames
- NewDumpHook
- WebSocket /ws client
- TestBatchInodeRows
- launch_probe.c
- install.sh
- mc_tag.h
- IntegrationSuite
- make test-integration (rootful Docker suite)
- ParseEventsFlag
- supersede_probe.c
- bpfLsmListed
- trace-app-libs.sh
- IntegrationSuite
- hashCache
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- dialLiveSession
- formatEntry
- mounts.go
- reloadOnce
- Guard
- render.go
- bottomBar
- newTestVetter
- filevault.go
- NewMonitor
- daemon: permanent protection with encryption at rest
- bufio.Reader
- find.go
- openRoot
- find_test.go
- tempGrantSetup
- Installation
- os.File
- graphify-refresh.sh
- runMonitor
- hashFileShared
- runNetworkGuard
- inodeChain
- launchKey
- main_test.go
- BinaryStat
- Vault
- IntegrationSuite
- ResolveLibraryClosure
- Limitations
- Troubleshooting
- .markProcess
- mountsUnitTest
- Daemon configuration
- edit-protected
- IntegrationSuite
- parseGuardEvent
- Compatibility
- How it works
- openBunTmp
- io.Writer
- pathCache
- SanitizeText
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .eventTypeAllowed
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- statInodeKey
- markReady
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
8. `IntegrationSuite` - 67 edges
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

## Communities (199 total, 31 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (22): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+14 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (103): TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse() (+95 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "GuardInodeKey"
Cohesion: 0.07
Nodes (29): refusedReplacements, TestNearestRootClaims(), appendInodeRow(), batchInodeRows(), inodeRowsOf(), nearestRootClaims(), newVerifyState(), supersedeUnnamed() (+21 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (83): TestSetSectionWhitelistPreservesGroupStructure(), updateCatalogConfig(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines() (+75 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.12
Nodes (45): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir() (+37 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (19): guardUnitTest, ReplacementCheck, commMatchesGuardedBinary(), eventMask(), guardModeKey(), addRecorder(), copySelf(), runTool() (+11 more)

### Community 9 - "mkdirs"
Cohesion: 0.21
Nodes (24): buildGlobReservations(), bitOf(), discordLibConfig(), hasChild(), mkdirs(), steamLibDirConfig(), TestBuildGlobReservations_BunRootsListed(), TestBuildGlobReservations_BunTmpdir() (+16 more)

### Community 10 - "BinaryEntry"
Cohesion: 0.08
Nodes (18): BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestWithChmodDropWriteNeedsReadOnly(), checkBtrfsKeys(), BinariesSummary(), canonicalBinaryPath(), canonicalPaths(), confinedEntry() (+10 more)

### Community 11 - "liveSession"
Cohesion: 0.14
Nodes (17): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries() (+9 more)

### Community 12 - "go_pkg_path_filepath"
Cohesion: 0.32
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 13 - "Config"
Cohesion: 0.10
Nodes (33): newBinaryVetter(), openBinaryVetter(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard() (+25 more)

### Community 14 - "install.go"
Cohesion: 0.06
Nodes (62): reloadDaemonForPasswordChange(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups() (+54 more)

### Community 15 - "uninstall_system.go"
Cohesion: 0.08
Nodes (32): configPaths(), pathCovered(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), cleanOrphanedFscrypt(), isInsidePath(), pickUsers() (+24 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.18
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (10): ownerPaths(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth.go"
Cohesion: 0.14
Nodes (28): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+20 more)

### Community 24 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (70): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+62 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.11
Nodes (11): NetBpfEvent, NetEvent, NetEventType, Cstr(), FormatAddr(), Ntohs(), ParseNetEventType(), NewNetworkMonitor() (+3 more)

### Community 27 - "highlight_test.go"
Cohesion: 0.15
Nodes (24): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+16 more)

### Community 28 - "LibraryClosure"
Cohesion: 0.14
Nodes (14): ldConfParser, defaultLibDirs(), elfInterp(), isScript(), ldConfInclude(), ldSoConfDirs(), LibraryClosure(), rootOwnedInode() (+6 more)

### Community 29 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 30 - "orphans_test.go"
Cohesion: 0.18
Nodes (19): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies() (+11 more)

### Community 31 - "serve.go"
Cohesion: 0.06
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.17
Nodes (25): fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod() (+17 more)

### Community 35 - "Ledger"
Cohesion: 0.15
Nodes (3): Ledger, Pending, Source

### Community 36 - "uninstall.go"
Cohesion: 0.05
Nodes (64): confirmMasterKeyOverwrite(), runGenKey(), TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), confirmRunPrereq(), deletePostBackups(), restoreBackups() (+56 more)

### Community 38 - "controlServer"
Cohesion: 0.18
Nodes (8): readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), controlServer, editControlSession, grantRequest

### Community 39 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 40 - "releaseAssets"
Cohesion: 0.18
Nodes (12): assetsFor(), checkUpdatePreconditions(), resolveAssets(), sanityCheckBinary(), TestAssetsFor(), TestResolveAssetsPicksArch(), TestSanityCheckBinaryRefusesOtherMachine(), TestValidateChannel() (+4 more)

### Community 41 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 44 - "newWatchPattern"
Cohesion: 0.18
Nodes (18): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+10 more)

### Community 45 - "update_test.go"
Cohesion: 0.10
Nodes (10): fetchReleases(), readInstalledVersion(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), TestNewerThanInstalled(), TestNewerThanStable(), TestParseStableVersion() (+2 more)

### Community 46 - "expandPlaceholders"
Cohesion: 0.15
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "prepareDaemonStart"
Cohesion: 0.10
Nodes (23): CheckBPFLSM(), configureDaemonLogging(), loadDaemonConfig(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck(), runLockdown() (+15 more)

### Community 49 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 50 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 51 - "codeedit.go"
Cohesion: 0.25
Nodes (17): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+9 more)

### Community 52 - "Read"
Cohesion: 0.09
Nodes (21): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+13 more)

### Community 53 - "watchGroup"
Cohesion: 0.11
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.11
Nodes (8): clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 55 - "TrustGuard"
Cohesion: 0.05
Nodes (21): tempGrant, tempMaskOp, tempRow, trustHook, trustRows, deleteInoKeys(), deleteResKeys(), setSupersedeTrust() (+13 more)

### Community 56 - "system.go"
Cohesion: 0.12
Nodes (26): dupFD(), installFile(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit(), mustCwd() (+18 more)

### Community 58 - "engine"
Cohesion: 0.05
Nodes (28): engine, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), pathWithin() (+20 more)

### Community 60 - "CandidateDir"
Cohesion: 0.17
Nodes (11): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestGlobReservationsLibTrusted(), TestParseGlobName() (+3 more)

### Community 61 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.15
Nodes (19): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+11 more)

### Community 63 - "GuardEvent"
Cohesion: 0.07
Nodes (27): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless() (+19 more)

### Community 64 - "netGuardModel"
Cohesion: 0.14
Nodes (8): sanitizeTerminalText(), sanitizeTerminalTexts(), clipToWidth(), formatResourceBar(), protoString(), TestSanitizeTerminalText(), netGuardEventMsg, netGuardModel

### Community 65 - "pinstate.go"
Cohesion: 0.12
Nodes (31): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState() (+23 more)

### Community 66 - "gui/session.go"
Cohesion: 0.05
Nodes (29): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+21 more)

### Community 67 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (9): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+1 more)

### Community 68 - "changelog.go"
Cohesion: 0.22
Nodes (19): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), showChangelog(), compareStableVersions(), confirmUpdate(), filterChannel() (+11 more)

### Community 69 - "buildTrustedSet"
Cohesion: 0.16
Nodes (16): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+8 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "KernelDev"
Cohesion: 0.11
Nodes (24): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+16 more)

### Community 73 - "time.Time"
Cohesion: 0.16
Nodes (9): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID, sweepMemo (+1 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "confgen.go"
Cohesion: 0.06
Nodes (54): logRefreshChange(), TestLogRefreshChange(), TestDiffMergePreservesExistingAndParses(), sectionsFromCandidates(), LibraryBlock, RefreshOptions, Section, SectionScan (+46 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 78 - "Highlighter"
Cohesion: 0.12
Nodes (7): span, diffModel, editorModel, Highlighter, lexerFor(), shebangLexer(), changelogModel

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 82 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 83 - "time.Duration"
Cohesion: 0.16
Nodes (4): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), debouncer, IntegrationSuite

### Community 84 - "Pulse"
Cohesion: 0.28
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 85 - "downloadAndVerify"
Cohesion: 0.31
Nodes (7): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), updateHTTPClient(), verifyRelease()

### Community 86 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 87 - "fakeGuardRepo"
Cohesion: 0.05
Nodes (22): addErrKind, classifyAddErr(), walkEntries(), walkLiveEntries(), SupersededBinary, TempBinary, TemporaryGrant, TempRule (+14 more)

### Community 88 - "unit.go"
Cohesion: 0.19
Nodes (13): TestDaemonUnitWithMetadataOutput(), DaemonUnitWithMetadataOutput(), ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate() (+5 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "catalogRefresher"
Cohesion: 0.15
Nodes (11): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps() (+3 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.08
Nodes (40): ServeConfig, UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), newDaemonModel(), runDaemon(), runDaemonUI() (+32 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.11
Nodes (8): foreignPinnedLink(), FirstOnBtrfs(), OnBtrfs(), eventTypeKey(), mcNameKey(), statInodeKey(), hook, NetGuard

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 100 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.17
Nodes (12): placedWalk, binaryStatOf(), confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd() (+4 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "NetGuardEvent"
Cohesion: 0.11
Nodes (17): NetGuardEvent, NetworkGuardRepository, NewMonitorUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo() (+9 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "EventType"
Cohesion: 0.10
Nodes (12): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), Run() (+4 more)

### Community 111 - "sync.Mutex"
Cohesion: 0.21
Nodes (4): newControlManager(), startControlManager(), configEditor, controlManager

### Community 112 - "testKeyPair"
Cohesion: 0.24
Nodes (9): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+1 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - ".PathsFor"
Cohesion: 0.17
Nodes (13): systemPatterns(), systemRule(), Discover(), DiscoverSystem(), BinaryRule, TestDiscoverOnlyExisting(), TestDiscoverSystemOnlyAbsolute(), TestPathsFor() (+5 more)

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 121 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "TestBatchInodeRows"
Cohesion: 0.50
Nodes (5): inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), TestCountInodeRows()

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

### Community 129 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "hashCache"
Cohesion: 0.18
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "dialLiveSession"
Cohesion: 0.16
Nodes (15): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+7 more)

### Community 143 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 144 - "mounts.go"
Cohesion: 0.26
Nodes (7): joinMount(), parseKernelDev(), parseMountinfo(), underRoot(), unescapeMountinfo(), mountEntry, mountTable

### Community 145 - "reloadOnce"
Cohesion: 0.13
Nodes (19): backingDeviceUnion(), buildConcurrency(), buildGuards(), drainEphemeral(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon(), unlockPendingGroupRoots() (+11 more)

### Community 146 - "Guard"
Cohesion: 0.07
Nodes (10): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, exeRow, rootHandle, SharedPinDegraded() (+2 more)

### Community 147 - "render.go"
Cohesion: 0.13
Nodes (13): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 148 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 149 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 150 - "filevault.go"
Cohesion: 0.06
Nodes (48): deprovisionKind, migrationTarget, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), isFileVaultCiphertext() (+40 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "bufio.Reader"
Cohesion: 0.20
Nodes (11): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest(), TestReadGrantRequestForward() (+3 more)

### Community 154 - "find.go"
Cohesion: 0.20
Nodes (7): findButton, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 155 - "openRoot"
Cohesion: 0.22
Nodes (7): openRoot(), mountinfoDev(), SuperblockDev(), SuperblockDevPath(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs(), progressReader

### Community 156 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 157 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.06
Nodes (47): VettedInode, Multicall, MulticallKind, ExeKey(), IsMulticallRefusal(), multicallFamily(), MulticallRefusal(), RefuseMulticallLines() (+39 more)

### Community 161 - "runMonitor"
Cohesion: 0.13
Nodes (18): BuildEBPFTargets(), CheckEBPF(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair() (+10 more)

### Community 163 - "hashFileShared"
Cohesion: 0.29
Nodes (8): HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled(), unchangedSince()

### Community 164 - "runNetworkGuard"
Cohesion: 0.16
Nodes (18): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), NetEventTypes(), DiscoverInfraBinaries() (+10 more)

### Community 165 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 167 - "main_test.go"
Cohesion: 0.18
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 168 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 171 - "ResolveLibraryClosure"
Cohesion: 0.18
Nodes (11): hostView(), mountVouchesOwnership(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch() (+3 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - ".markProcess"
Cohesion: 0.22
Nodes (7): reasonText(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), TestReasonText(), procRec

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 183 - "io.Writer"
Cohesion: 0.43
Nodes (7): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), candidate

### Community 184 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 185 - "SanitizeText"
Cohesion: 0.10
Nodes (25): multicallAdmissible(), buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), binaryVetter, vetDecision, AdmissionCheck, deferredBinary (+17 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

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

### Community 195 - "markReady"
Cohesion: 0.50
Nodes (4): markReady(), notifySystemdReady(), InodeOccupancy(), WarnInodeOccupancy()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.12
Nodes (16): checkFreshDir(), openDirNoFollow(), renewDir(), TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), isBinaryContent(), listEntries() (+8 more)

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 452 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **31 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `newWatchPattern` to `Daemon configuration`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12223599063747763 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `newWatchPattern`, `README.md`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.031805929919137464 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `newWatchPattern` to `time.Time`, `catalogRefresher`, `go_pkg_path_filepath`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09056732740943267 - nodes in this community are weakly interconnected._