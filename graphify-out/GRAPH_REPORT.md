# Graph Report - app-listener  (2026-10-12)

## Corpus Check
- 352 files · ~1,554,367 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4594 nodes · 18207 edges · 196 communities (162 shown, 34 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1307 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6904a65b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- install.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- IntegrationSuite
- mayUpdate
- TestMonitorUseCaseLifecycle
- Monitor
- systemd.go
- Vault
- IntegrationSuite
- guard_trust.bpf.c
- confgen.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- ResolvePinBase
- reloadOnce
- guard.bpf.c
- NetworkMonitor
- find_test.go
- Ledger
- buildOneGuard
- monitorUnitTest
- serve.go
- engine
- IntegrationSuite
- fileedit_model_test.go
- IntegrationSuite
- Run
- uninstall_system.go
- controlServer
- bpfstats/main.go
- github.com/charmbracelet/bubbles/textarea.Model
- migrate.go
- absPath
- liveSession
- IntegrationSuite
- update.go
- highlight_test.go
- .syncSetsLocked
- TrustGuard
- tempRow
- inode_dev
- IntegrationSuite
- applyConfigLine
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- codeedit.go
- sync.Mutex
- filevault_test.go
- parseConfig
- update_test.go
- runNetworkGuard
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- Candidate
- CandidateDir
- filevault.go
- time.Duration
- fcntl
- copyRegularAt
- Config
- SanitizeText
- btrfsLayoutFrom
- New
- github.com/spf13/cobra.Command
- runNetworkMonitor
- process_vm_readv.c
- OpenSystemPlaced
- serveWebSocket
- Read
- Vault
- MountVouch
- safeio.go
- render.go
- time.Time
- .GrantTemporaryAccess
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- auditGroupCoverage
- KernelDev
- IntegrationSuite
- Well-known bypass classes
- model
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- supersedeMaps
- Build & sign release assets (reusable workflow)
- runMonitor
- swappedParent
- net_tester/main.go
- NewEventFanout
- app-listener Terminal Demo Recording
- IntegrationSuite
- watchGroup
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- catalogRefresher
- buildTrustedSet
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
- fakeDaemon
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- writeWithin
- ResolveLibraryClosure
- os.File
- preload_wait.c
- Guard
- integrationtests/guard_test.go
- launchKey
- fscrypt/orphans.go
- find.go
- is_event_type_allowed
- daemon: permanent protection with encryption at rest
- newTestVetter
- GuardInodeKey
- bottomBar
- .pathInfo
- Pulse
- Installation
- tempGrantSetup
- graphify-refresh.sh
- IntegrationSuite
- .degradeWatchLoop
- github.com/cilium/ebpf/link.Link
- NewMonitor
- GuardEvent
- newerThanStable
- vouchedDevs
- Vault
- NetGuard
- parseGuardEvent
- Limitations
- Troubleshooting
- downloadFile
- launch_rule_at
- Daemon configuration
- edit-protected
- parseResize
- formatEntry
- Compatibility
- How it works
- harnessSuite
- TestMain
- bpfLsmListed
- .runEditorHarness
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS

## God Nodes (most connected - your core abstractions)
1. `Guard` - 95 edges
2. `Config` - 93 edges
3. `GuardInodeKey` - 91 edges
4. `Load()` - 85 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
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

## Communities (196 total, 34 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (24): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+16 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (125): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+117 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "install.go"
Cohesion: 0.12
Nodes (22): promptEditPassword(), applyLiveRefresh(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), encryptDirectories(), ensureInstalledBinary(), parseMaintenanceFlags(), prepareInstallation() (+14 more)

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
Cohesion: 0.10
Nodes (52): trustManager, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup() (+44 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (22): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf() (+14 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "IntegrationSuite"
Cohesion: 0.07
Nodes (43): ParseEventsFlag(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+35 more)

### Community 11 - "mayUpdate"
Cohesion: 0.20
Nodes (11): hardLinked(), mayUpdate(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), TestMayUpdate_HardLinkWaivedForUserWritableLibBinary(), TestPlanUpdaters_ScopedToOwnResources() (+3 more)

### Community 12 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.21
Nodes (10): NewGuardUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkMonitorUseCaseLifecycle() (+2 more)

### Community 14 - "systemd.go"
Cohesion: 0.05
Nodes (73): reloadDaemonForPasswordChange(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), buildBinaryIfNeeded(), ensureMasterKey(), installBinaryOnly(), runUpdateCatalogOnly(), collectFilesystemPrereqs() (+65 more)

### Community 15 - "Vault"
Cohesion: 0.12
Nodes (19): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+11 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.12
Nodes (3): IntegrationSuite, guardDeltaEvents(), le64HexKey()

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "confgen.go"
Cohesion: 0.06
Nodes (53): logRefreshChange(), TestLogRefreshChange(), appendSectionsAndEdit(), libraryBlocksFromCandidates(), LibraryBlock, RefreshOptions, Section, findSectionEnd() (+45 more)

### Community 19 - "fileEditModel"
Cohesion: 0.11
Nodes (10): clipLabel(), humanSize(), isBinaryContent(), sortNodes(), TestIsBinaryContent(), chownUser, fileEditMode, fileEditModel (+2 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (16): backingDeviceUnion(), libBinaries(), planUpdaters(), resourceBinaries(), resourceLinks(), TestPlanUpdaters_LinkOwnsItsResourceBit(), Resource, BackingDevices() (+8 more)

### Community 23 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 24 - "reloadOnce"
Cohesion: 0.11
Nodes (25): CheckEBPF(), buildConcurrency(), buildGuards(), lockOneRoot(), lockRootRecovering(), makeReloadHandler(), newPinGeneration(), prepareDaemonStart() (+17 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (51): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), evict_inode_from_guard(), get_dentry_from_path() (+43 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.13
Nodes (9): NetBpfEvent, NetEvent, FormatAddr(), Ntohs(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey() (+1 more)

### Community 27 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 28 - "Ledger"
Cohesion: 0.10
Nodes (16): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), EnsurePlaceholders(), Ledger (+8 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.10
Nodes (27): buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions(), newSelfGuards(), selfGuards, AdmissionCheck (+19 more)

### Community 31 - "serve.go"
Cohesion: 0.12
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

### Community 32 - "engine"
Cohesion: 0.09
Nodes (8): TestTieRank(), engine, pathWithin(), TestPathWithinMatchesPrefixRule(), tieRank(), deleteInoKeys(), intersectAllows(), TestIntersectAllows()

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.16
Nodes (25): fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod() (+17 more)

### Community 35 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 36 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 37 - "uninstall_system.go"
Cohesion: 0.09
Nodes (35): runMaintenanceMode(), deletePostBackups(), restoreBackups(), pickUsers(), offerBackupCleanup(), runUninstall(), detectSSHAgentUnits(), isInstallerSSHAgentUnit() (+27 more)

### Community 38 - "controlServer"
Cohesion: 0.22
Nodes (5): controlServer, readClientLines(), streamEvents(), editControlSession, grantRequest

### Community 39 - "bpfstats/main.go"
Cohesion: 0.13
Nodes (23): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection() (+15 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.17
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "update.go"
Cohesion: 0.11
Nodes (32): changelogText(), showChangelog(), TestChangelogText(), assetsFor(), checkUpdatePreconditions(), confirmUpdate(), downloadAndVerify(), downloadReleaseFiles() (+24 more)

### Community 46 - "highlight_test.go"
Cohesion: 0.15
Nodes (16): Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens(), TestGutterKeepsTextColumnAligned() (+8 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "TrustGuard"
Cohesion: 0.11
Nodes (8): trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 50 - "tempRow"
Cohesion: 0.18
Nodes (6): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust()

### Community 51 - "inode_dev"
Cohesion: 0.12
Nodes (23): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), code_suspect(), current_is_inspector() (+15 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 53 - "applyConfigLine"
Cohesion: 0.14
Nodes (23): applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary(), applyLibDir(), applyLibraryBlockDirective() (+15 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.11
Nodes (7): clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 56 - "codeedit.go"
Cohesion: 0.25
Nodes (17): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+9 more)

### Community 57 - "sync.Mutex"
Cohesion: 0.17
Nodes (6): controlManager, controlManager, newControlManager(), startControlManager(), startControlServer(), configEditor

### Community 58 - "filevault_test.go"
Cohesion: 0.18
Nodes (26): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsEncryptedDispatchesToFileVaultForRegularFiles() (+18 more)

### Community 59 - "parseConfig"
Cohesion: 0.25
Nodes (7): parseConfig(), entryOf(), BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype()

### Community 60 - "update_test.go"
Cohesion: 0.11
Nodes (10): rsaPublicKeyPEM(), signedChecksum(), TestAssetsFor(), TestFilterChannel(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestReadInstalledVersion() (+2 more)

### Community 61 - "runNetworkGuard"
Cohesion: 0.09
Nodes (19): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode (+11 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 64 - "errno"
Cohesion: 0.28
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.10
Nodes (33): configureDaemonLogging(), loadDaemonConfig(), resolveConfigPath(), runLockdown(), setConfinedHomes(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive() (+25 more)

### Community 66 - "Candidate"
Cohesion: 0.07
Nodes (36): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), collectDiffAdditions(), configPaths(), pathCovered(), runDiffCatalog(), TestPathCovered() (+28 more)

### Community 67 - "CandidateDir"
Cohesion: 0.06
Nodes (32): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, compileGlobNames(), globText(), ParseGlobName() (+24 more)

### Community 68 - "filevault.go"
Cohesion: 0.23
Nodes (14): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+6 more)

### Community 69 - "time.Duration"
Cohesion: 0.16
Nodes (9): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer, New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod() (+1 more)

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "Config"
Cohesion: 0.22
Nodes (14): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager, startTrustGuard(), trustStartupError(), selectAndEditConfig(), editConfig(), runConfigEditor() (+6 more)

### Community 73 - "SanitizeText"
Cohesion: 0.19
Nodes (7): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), binaryVetter, vetDecision, SanitizeText()

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 75 - "New"
Cohesion: 0.17
Nodes (13): diffModel, editorModel, moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys() (+5 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.09
Nodes (30): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+22 more)

### Community 78 - "runNetworkMonitor"
Cohesion: 0.15
Nodes (9): runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel(), NetworkMonitorUseCase (+1 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "OpenSystemPlaced"
Cohesion: 0.13
Nodes (18): placedWalk, fileKey(), TestSystemFileRefusesAnotherInode(), TestSystemFileRefusesUserPlacedName(), confineBelow(), inUserHome(), GlobSystemPlaced(), OpenSystemPlaced() (+10 more)

### Community 81 - "serveWebSocket"
Cohesion: 0.20
Nodes (8): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 82 - "Read"
Cohesion: 0.09
Nodes (22): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+14 more)

### Community 83 - "Vault"
Cohesion: 0.23
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 84 - "MountVouch"
Cohesion: 0.07
Nodes (18): GlobChild, statFunc, syncU32Map(), loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows(), TrustGuard, resolveBits() (+10 more)

### Community 85 - "safeio.go"
Cohesion: 0.07
Nodes (40): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+32 more)

### Community 86 - "render.go"
Cohesion: 0.13
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "time.Time"
Cohesion: 0.09
Nodes (19): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+11 more)

### Community 88 - ".GrantTemporaryAccess"
Cohesion: 0.14
Nodes (12): TemporaryGrant, TempRule, applyTempGrants(), daemonUseCase, TemporaryAccess, logTempGrant(), newExeTap(), planTempGrants() (+4 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.12
Nodes (17): ldConfParser, defaultLibDirs(), elfInterp(), isScript(), ldConfInclude(), ldSoConfDirs(), LibraryClosure(), resolveSoname() (+9 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.09
Nodes (38): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), markReady(), newDaemonModel(), notifySystemdReady() (+30 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.11
Nodes (12): NetEventType, FirstOnBtrfs(), OnBtrfs(), TestSuperblockDevMatchesStatOffBtrfs(), NetEventTypes(), ParseNetEventType(), eventsetSummary(), eventTypeKey() (+4 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 100 - "KernelDev"
Cohesion: 0.07
Nodes (38): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+30 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "model"
Cohesion: 0.22
Nodes (7): formatType(), listenForEvents(), syncViewport(), tickStats(), eventLine, eventMsg, model

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "supersedeMaps"
Cohesion: 0.15
Nodes (11): checkBtrfsKeys(), ensureBtrfsLayout(), dropSupersededLocked(), Guard, PruneSuperseded(), supersededMark(), supersedeMaps(), supersedePruneLoop() (+3 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.11
Nodes (19): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+11 more)

### Community 110 - "swappedParent"
Cohesion: 0.18
Nodes (15): AtomicWriteAt(), openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent() (+7 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NewEventFanout"
Cohesion: 0.29
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout[T]

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "watchGroup"
Cohesion: 0.26
Nodes (12): watchGroup, addResource(), deferPendingWatchPath(), finalizeGroup(), findResource(), isInsidePath(), materializeLibDirs(), materializeWatchGroup() (+4 more)

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "catalogRefresher"
Cohesion: 0.13
Nodes (11): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), grantNewBinaries(), catalogRefresher, Line() (+3 more)

### Community 121 - "buildTrustedSet"
Cohesion: 0.15
Nodes (17): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+9 more)

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
Cohesion: 0.23
Nodes (6): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "EventType"
Cohesion: 0.07
Nodes (16): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock() (+8 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 143 - "ResolveLibraryClosure"
Cohesion: 0.15
Nodes (13): fileExists(), hostView(), LdSoPreloadPaths(), mountVouchesOwnership(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace() (+5 more)

### Community 144 - "os.File"
Cohesion: 0.06
Nodes (55): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+47 more)

### Community 146 - "Guard"
Cohesion: 0.05
Nodes (32): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, deferredBinary, rootHandle, sharedHashEntry (+24 more)

### Community 147 - "integrationtests/guard_test.go"
Cohesion: 0.36
Nodes (6): attrOpCase, eventsForPath(), guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 148 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 149 - "fscrypt/orphans.go"
Cohesion: 0.22
Nodes (20): cleanOrphanedFscrypt(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 150 - "find.go"
Cohesion: 0.27
Nodes (7): findButton, span, findAll(), lowerRunes(), newFindInput(), runesEqual(), TestFindAllSmartCase()

### Community 151 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "newTestVetter"
Cohesion: 0.42
Nodes (9): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSwapWithinGenerationLosesRights() (+1 more)

### Community 154 - "GuardInodeKey"
Cohesion: 0.07
Nodes (23): inodeBatchFn, refusedReplacements, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), engine (+15 more)

### Community 155 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 156 - ".pathInfo"
Cohesion: 0.48
Nodes (3): evalSymlinksOrEmpty(), pathCache, pathCacheEntry

### Community 157 - "Pulse"
Cohesion: 0.28
Nodes (3): Pulse, Sub, EventFanout

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 163 - "github.com/cilium/ebpf/link.Link"
Cohesion: 0.29
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 165 - "GuardEvent"
Cohesion: 0.07
Nodes (23): selfProtectSpecs(), displayEntries(), holdForward(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless() (+15 more)

### Community 166 - "newerThanStable"
Cohesion: 0.33
Nodes (7): compareStableVersions(), newerThanStable(), parseStableVersion(), parseVerPart(), TestNewerThanStable(), TestParseStableVersion(), stableVersion

### Community 167 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 170 - "parseGuardEvent"
Cohesion: 0.22
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "downloadFile"
Cohesion: 0.33
Nodes (4): downloadFile(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), progressReader

### Community 175 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 179 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 183 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 186 - "Binary trust"
Cohesion: 0.33
Nodes (6): Binary ledger, Binary replacement, Binary trust, Launch scan, Multicall binaries, trust-binaries

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

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.18
Nodes (9): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), applyNewFileMeta(), defaultNewFileMode(), listEntries(), TestDefaultNewFileMode(), writeFileInPlace(), dirEntry (+1 more)

## Knowledge Gaps
- **145 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+140 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 480 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **34 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `GuardEvent`, `fakeDaemon`, `os.File`, `Pulse`, `time.Time`, `reloadOnce`, `sync.Mutex`, `GuardInodeKey`, `DaemonEvent`, `buildOneGuard`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **Why does `engine` connect `engine` to `daemon/daemon.go`, `github.com/cilium/ebpf/link.Link`, `reloadOnce`, `sync.Mutex`, `GuardInodeKey`, `NetGuard`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _145 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11067377340440071 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02374507874015748 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._