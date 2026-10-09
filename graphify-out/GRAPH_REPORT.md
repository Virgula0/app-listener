# Graph Report - app-listener  (2026-10-09)

## Corpus Check
- 345 files · ~1,549,435 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4536 nodes · 17925 edges · 199 communities (164 shown, 35 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1294 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a1cce841`
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
- Hash
- planUpdaters
- codeedit.go
- uninstall_system.go
- install.go
- Vault
- IntegrationSuite
- guard_trust.bpf.c
- confgen.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- RawTarget
- runNetworkGuard
- guard.bpf.c
- app.go
- github.com/spf13/cobra.Command
- reloadOnce
- buildOneGuard
- Monitor
- serve.go
- GuardInodeKey
- time.Duration
- fileedit_model_test.go
- DaemonEvent
- watchGroup
- New
- GuardEvent
- bpfstats/main.go
- github.com/charmbracelet/bubbles/textarea.Model
- NetEventType
- absPath
- liveSession
- IntegrationSuite
- fscrypt/orphans.go
- monitorUnitTest
- .syncSetsLocked
- guardModel
- tempRow
- __always_inline
- IntegrationSuite
- TrustGuard
- github.com/charmbracelet/bubbletea.Model
- IntegrationSuite
- SanitizeText
- sync.Mutex
- github.com/stretchr/testify/suite.Suite
- EventType
- vaultFS
- netGuardModel
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- netModel
- highlight.go
- LibraryClosure
- catalogRefresher
- fcntl
- copyRegularAt
- dataRowRenderer
- MountVouch
- btrfsLayoutFrom
- Ledger
- RefreshSection
- runMonitor
- process_vm_readv.c
- model
- backups.go
- update_test.go
- expandPlaceholders
- NetworkMonitor
- safeio.go
- render.go
- watchSet
- TempBinary
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- KernelDev
- fscrypt.go
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- serveWebSocket
- sanitizeTerminalText
- IntegrationSuite
- Well-known bypass classes
- Vault
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- NetGuard
- Build & sign release assets (reusable workflow)
- time.Time
- runGuard
- net_tester/main.go
- NetworkGuardUseCase
- app-listener Terminal Demo Recording
- is_event_type_allowed
- .step
- fileedit_symlink_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- auditGroupCoverage
- trustManager
- WebSocket /ws client
- NewDumpHook
- launch_probe.c
- install.sh
- mc_tag.h
- ListUsers
- make test-integration (rootful Docker suite)
- NewEventFanout
- supersede_probe.c
- IntegrationSuite
- trace-app-libs.sh
- IntegrationSuite
- Config
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- bottomBar
- find_test.go
- os.File
- Guard
- runEditProtected
- update.go
- configedit_test.go
- inode_dev
- writeWithin
- daemon: permanent protection with encryption at rest
- IntegrationSuite
- CandidateDir
- tempGrantSetup
- launchKey
- runNetworkMonitor
- Installation
- lockRootRecovering
- graphify-refresh.sh
- downloadFile
- runDaemon
- newChangelogModel
- Candidate
- Guard
- launch_rule_at
- IntegrationSuite
- .pathInfo
- AddServeFlags
- ParseEventsFlag
- Limitations
- Troubleshooting
- Verify
- fakeDaemon
- Daemon configuration
- edit-protected
- eventUnitTest
- loadDaemonConfig
- Compatibility
- How it works
- rootHandle
- Cstr
- bpfLsmListed
- parseResize
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- formatEntry
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- Vault
- TestMain
- .runEditorHarness
- monitor: see what touches a path
- newPathCache

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

## Communities (199 total, 35 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (28): atomicWriteAt(), writeAndSync(), Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent (+20 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (121): plantBun(), TestOpenBunTmp_KeepsVettedDir(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait() (+113 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "filevault.go"
Cohesion: 0.14
Nodes (28): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, isFileVaultCiphertext(), looksLikeFileVaultRecord() (+20 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (69): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+61 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.11
Nodes (45): trustManager, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), addKeysToAgent(), rcTarget, BunTmpDir(), catalogEntryMatch() (+37 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf() (+13 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "Hash"
Cohesion: 0.19
Nodes (19): Hash(), HashFileExists(), LoadHashFile(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle(), TestRemoveHashFileInPlace(), TestValidatePassword() (+11 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "codeedit.go"
Cohesion: 0.10
Nodes (41): nextConfig(), class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+33 more)

### Community 14 - "install.go"
Cohesion: 0.06
Nodes (64): reloadDaemonForPasswordChange(), askBunTmpdirUsers(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), applyLiveRefresh() (+56 more)

### Community 15 - "Vault"
Cohesion: 0.14
Nodes (17): deprovisionKind, isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn() (+9 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "confgen.go"
Cohesion: 0.09
Nodes (41): TestSetSectionWhitelistPreservesGroupStructure(), LibraryBlock, Section, ConfSafePath(), findSectionEnd(), findSectionStart(), firstContentLine(), GenerateConf() (+33 more)

### Community 19 - "fileEditModel"
Cohesion: 0.09
Nodes (13): clipLabel(), fileModeToUnixMode(), humanSize(), clamp(), sortNodes(), listenForGuardEvents(), tickGuardStats(), chownUser (+5 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (11): backingDeviceUnion(), Resource, BackingDevices(), GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots() (+3 more)

### Community 23 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 24 - "runNetworkGuard"
Cohesion: 0.30
Nodes (11): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), eventsetSummary(), Mode (+3 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.14
Nodes (37): add_inode_to_guard(), chmod_only_drops_write(), discover_guarded_parent(), evict_inode_from_guard(), fill_path(), get_dentry_from_path(), get_inode_from_path(), guard_file_truncate() (+29 more)

### Community 26 - "app.go"
Cohesion: 0.10
Nodes (12): columnState, dataRow, guiModel, headerWidget, findTopLevelWindow(), floatWindow(), internAtom(), newColumnState() (+4 more)

### Community 27 - "github.com/spf13/cobra.Command"
Cohesion: 0.26
Nodes (12): credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink() (+4 more)

### Community 28 - "reloadOnce"
Cohesion: 0.15
Nodes (17): buildConcurrency(), buildGuards(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded(), startCatalogRefresh(), startGuardedDaemon(), startGuardedDaemonAbortable() (+9 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.11
Nodes (25): buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions(), selfGuards, AdmissionCheck, VettedInode (+17 more)

### Community 31 - "serve.go"
Cohesion: 0.12
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

### Community 32 - "GuardInodeKey"
Cohesion: 0.05
Nodes (23): refusedReplacements, TestTieRank(), engine, pathWithin(), tieRank(), deleteInoKeys(), engine, Inspector (+15 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.17
Nodes (24): isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod() (+16 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.18
Nodes (19): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+11 more)

### Community 36 - "watchGroup"
Cohesion: 0.15
Nodes (25): watchGroup, applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary(), applyLibDir() (+17 more)

### Community 37 - "New"
Cohesion: 0.12
Nodes (27): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+19 more)

### Community 38 - "GuardEvent"
Cohesion: 0.08
Nodes (15): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), bpfGuardEvent, commMatchesGuardedBinary(), fsGateLabel() (+7 more)

### Community 39 - "bpfstats/main.go"
Cohesion: 0.17
Nodes (18): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), InertProgram(), KernelHasFunc(), TestKernelHasFunc(), trustSpec() (+10 more)

### Community 40 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (10): findButton, span, reposition(), findAll(), Finder, lowerRunes(), newFindInput(), overlay() (+2 more)

### Community 41 - "NetEventType"
Cohesion: 0.12
Nodes (9): NetEvent, NetEventType, ParseNetEventType(), NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase(), newFakeNetworkMonitorRepo(), TestNetworkMonitorUseCaseLifecycle() (+1 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.08
Nodes (15): guardBinaryFlag(), infraContainerPaths(), netMonitorEvent, IntegrationSuite, IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount() (+7 more)

### Community 45 - "fscrypt/orphans.go"
Cohesion: 0.18
Nodes (21): cleanOrphanedFscrypt(), modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors() (+13 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "guardModel"
Cohesion: 0.12
Nodes (17): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+9 more)

### Community 50 - "tempRow"
Cohesion: 0.29
Nodes (4): tempGrant, tempMaskOp, tempRow, liveTrust()

### Community 51 - "__always_inline"
Cohesion: 0.15
Nodes (24): exe_is_superseded(), exe_refused(), check_and_emit_args(), code_suspect(), count_degrade(), current_is_inspector(), emit_process_denial(), emit_with() (+16 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 53 - "TrustGuard"
Cohesion: 0.09
Nodes (12): statFunc, trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust(), cStr(), resolveBits() (+4 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.09
Nodes (8): editorModel, renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel, changelogModel

### Community 56 - "SanitizeText"
Cohesion: 0.15
Nodes (11): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), runHeadless(), runHeadless(), binaryVetter, vetDecision (+3 more)

### Community 57 - "sync.Mutex"
Cohesion: 0.08
Nodes (23): controlManager, controlServer, readConfigPut(), TestReadGrantRequestConfig(), controlManager, controlServer, newControlManager(), readAuthRequest() (+15 more)

### Community 59 - "EventType"
Cohesion: 0.12
Nodes (12): BpfEvent, EventType, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock(), Run() (+4 more)

### Community 60 - "vaultFS"
Cohesion: 0.12
Nodes (15): writeTempWithMetadata(), readClassifiedFile(), stampCopyMetadata(), TestCleanupStalePins(), TestRuntimeClassMarkers(), OpenRegularNoFollow(), applyNewFileMeta(), defaultNewFileMode() (+7 more)

### Community 61 - "netGuardModel"
Cohesion: 0.24
Nodes (5): NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 64 - "errno"
Cohesion: 0.21
Nodes (5): put(), wait_then_read(), dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (23): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+15 more)

### Community 66 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 67 - "highlight.go"
Cohesion: 0.19
Nodes (5): diffModel, Highlighter, lexerFor(), shebangLexer(), tokenClassOf()

### Community 68 - "LibraryClosure"
Cohesion: 0.08
Nodes (28): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+20 more)

### Community 69 - "catalogRefresher"
Cohesion: 0.18
Nodes (5): drain(), newCatalogRefresher(), readEvents(), catalogRefresher, debouncer

### Community 70 - "fcntl"
Cohesion: 0.13
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 73 - "MountVouch"
Cohesion: 0.09
Nodes (15): dropSupersededLocked(), Guard, PruneSuperseded(), supersededMark(), supersedeMaps(), supersedePruneLoop(), warnLostStamps(), syncU32Map() (+7 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "Ledger"
Cohesion: 0.08
Nodes (25): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSwapWithinGenerationLosesRights() (+17 more)

### Community 76 - "RefreshSection"
Cohesion: 0.15
Nodes (16): logRefreshChange(), TestLogRefreshChange(), RefreshOptions, SectionScan, admitNew(), SectionChange, keepExisting(), LiveEmptyWhitelistRejected() (+8 more)

### Community 78 - "runMonitor"
Cohesion: 0.16
Nodes (9): MakeDisplayPaths(), runHeadless(), runMonitor(), runTUI(), MonitorRepository, NewModel(), MonitorUseCase, NewMonitorUseCase() (+1 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "model"
Cohesion: 0.26
Nodes (6): formatType(), listenForEvents(), tickStats(), eventLine, eventMsg, model

### Community 81 - "backups.go"
Cohesion: 0.11
Nodes (28): deletePostBackups(), restoreBackups(), pickUsers(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir() (+20 more)

### Community 82 - "update_test.go"
Cohesion: 0.08
Nodes (19): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), testKeyPair() (+11 more)

### Community 83 - "expandPlaceholders"
Cohesion: 0.11
Nodes (15): symlinkTargetGlob(), confSafeMatch(), expandPlaceholders(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), CandidateDir, boundedGlob() (+7 more)

### Community 84 - "NetworkMonitor"
Cohesion: 0.16
Nodes (6): NetEventTypes(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 85 - "safeio.go"
Cohesion: 0.07
Nodes (37): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), TestDaemonUnitWithMetadataOutput() (+29 more)

### Community 86 - "render.go"
Cohesion: 0.14
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "watchSet"
Cohesion: 0.16
Nodes (12): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern (+4 more)

### Community 88 - "TempBinary"
Cohesion: 0.08
Nodes (25): BinaryStat, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow() (+17 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "KernelDev"
Cohesion: 0.06
Nodes (44): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+36 more)

### Community 92 - "fscrypt.go"
Cohesion: 0.07
Nodes (37): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), confirmRunPrereq(), resolveFilesystemPrereqs(), checkKeyLen(), classifySetupError(), classifySupportError() (+29 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.13
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

### Community 99 - "serveWebSocket"
Cohesion: 0.26
Nodes (8): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 100 - "sanitizeTerminalText"
Cohesion: 0.17
Nodes (8): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), TestSanitizeTerminalText()

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "Vault"
Cohesion: 0.35
Nodes (3): migrationTarget, classifyMigrationTarget(), Vault

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

### Community 109 - "time.Time"
Cohesion: 0.20
Nodes (13): newWatchSet(), newGateLogLimiter(), gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters() (+5 more)

### Community 110 - "runGuard"
Cohesion: 0.18
Nodes (11): selfProtectSpecs(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), selfProtectSpec, Mode, NewGuardModel() (+3 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NetworkGuardUseCase"
Cohesion: 0.17
Nodes (9): NetworkGuardRepository, NetworkGuardUseCase, NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkGuardUseCaseStartError() (+1 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 115 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 116 - "fileedit_symlink_test.go"
Cohesion: 0.24
Nodes (15): AtomicWriteAt(), openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent() (+7 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 121 - "trustManager"
Cohesion: 0.13
Nodes (19): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), trustManager (+11 more)

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

### Community 127 - "ListUsers"
Cohesion: 0.08
Nodes (28): configPaths(), pathCovered(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), isInsidePath(), askSSHAgentUsers(), askDecrypt() (+20 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "NewEventFanout"
Cohesion: 0.22
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "Config"
Cohesion: 0.13
Nodes (23): sameConfig(), swapFile(), bunEntryConfigured(), bunLaunchersFor(), grantNewBinaries(), Line(), Lines(), TestLinesKeyOnLinkAndIncludePending() (+15 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 143 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 144 - "os.File"
Cohesion: 0.06
Nodes (48): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp (+40 more)

### Community 146 - "Guard"
Cohesion: 0.06
Nodes (29): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, deferredBinary, sharedHashEntry, BinaryEntry (+21 more)

### Community 147 - "runEditProtected"
Cohesion: 0.16
Nodes (15): LiveModeAvailable(), runEditProtected(), chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries(), holdForward() (+7 more)

### Community 148 - "update.go"
Cohesion: 0.12
Nodes (33): changelogText(), showChangelog(), TestChangelogText(), applyUpdate(), assetsFor(), checkUpdatePreconditions(), compareStableVersions(), confirmUpdate() (+25 more)

### Community 149 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 150 - "inode_dev"
Cohesion: 0.16
Nodes (12): ctx_ptr(), inode_dev(), sb_dev(), btrfs_copy_gate(), guard_bprm_committed(), guard_exec_applet(), guard_file_ioctl(), guard_file_ioctl_compat() (+4 more)

### Community 151 - "writeWithin"
Cohesion: 0.29
Nodes (7): descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeWithin()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 154 - "CandidateDir"
Cohesion: 0.12
Nodes (16): entryWriters(), reserveChain(), underResource(), globBuilder, GlobChild, compileGlobNames(), globText(), GlobReservations (+8 more)

### Community 155 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 156 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 157 - "runNetworkMonitor"
Cohesion: 0.32
Nodes (6): ParseNetEventsFlag(), runNetworkMonitor(), runTUI(), NewNetModel(), assertViewFits(), TestViewsFitReportedTerminalSize()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "lockRootRecovering"
Cohesion: 0.29
Nodes (8): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), PinPrefix(), sanitizeGen(), SharedPinPrefix(), TestPinPrefixDistinctAndStable(), TestPinPrefixShape()

### Community 161 - "downloadFile"
Cohesion: 0.33
Nodes (4): downloadFile(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), progressReader

### Community 162 - "runDaemon"
Cohesion: 0.09
Nodes (25): CheckBPFLSM(), CheckEBPF(), catchLifecycleSignals(), configureDaemonLogging(), confirmMasterKeyOverwrite(), newPinGeneration(), notifySystemdReady(), prepareDaemonStart() (+17 more)

### Community 163 - "newChangelogModel"
Cohesion: 0.53
Nodes (5): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes()

### Community 164 - "Candidate"
Cohesion: 0.16
Nodes (21): appendSectionsAndEdit(), collectDiffAdditions(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories() (+13 more)

### Community 166 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 168 - ".pathInfo"
Cohesion: 0.39
Nodes (3): evalSymlinksOrEmpty(), pathCache, pathCacheEntry

### Community 169 - "AddServeFlags"
Cohesion: 0.29
Nodes (7): AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init()

### Community 170 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "Verify"
Cohesion: 0.43
Nodes (7): OriginOf(), parseHash(), TestHashVerifyRoundTrip(), TestVerifyMalformed(), Verify(), Origin, parsedHash

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), SetUserHomes(), TestOpenSystemPlacedRefusesUserHome()

### Community 180 - "Compatibility"
Cohesion: 0.29
Nodes (7): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "rootHandle"
Cohesion: 0.12
Nodes (9): rootHandle, openRoot(), entryOf(), BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), SuperblockDev(), TestMountinfoDev() (+1 more)

### Community 183 - "Cstr"
Cohesion: 0.33
Nodes (4): NetBpfEvent, Cstr(), FormatAddr(), Ntohs()

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 185 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 186 - "Binary trust"
Cohesion: 0.33
Nodes (6): Binary ledger, Binary replacement, Binary trust, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 190 - "Development"
Cohesion: 0.40
Nodes (5): Demo GIFs, Development, Docker, Makefile targets, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 195 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **141 isolated node(s):** `Whitelisted programs that disclose on request`, `Interpreters and JITs`, `Root`, `Binaries you can write`, `Event masks are not confinement` (+136 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 473 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **35 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `GuardInodeKey`, `DaemonEvent`, `GuardEvent`, `time.Time`, `runGuard`, `fakeDaemon`, `os.File`, `rootHandle`, `sync.Mutex`, `reloadOnce`, `buildOneGuard`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Why does `Monitor` connect `Monitor` to `daemon/daemon.go`, `.pathInfo`, `NetworkMonitor`, `sync.Mutex`, `reloadOnce`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **What connects `Whitelisted programs that disclose on request`, `Interpreters and JITs`, `Root` to the rest of the system?**
  _141 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11383624067570314 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.024914765276685024 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._