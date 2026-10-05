# Graph Report - app-listener  (2026-10-05)

## Corpus Check
- 308 files · ~1,481,576 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4062 nodes · 16164 edges · 165 communities (136 shown, 29 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1168 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e35048a2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- update.go
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- mkdirs
- auth.go
- Guard
- Monitor
- uninstall.go
- install.go
- runDaemon
- IntegrationSuite
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- Candidate
- usecase/daemon_test.go
- Resource
- CandidateDir
- go_pkg_os
- guard.bpf.c
- dataRowRenderer
- NetEventType
- SanitizeText
- buildOneGuard
- buildTrustedSet
- serve.go
- engine
- os.File
- fileedit_model_test.go
- DaemonEvent
- Config
- runNetworkGuard
- GuardEvent
- KernelDev
- backups.go
- codeedit.go
- absPath
- New
- main_test.go
- system.go
- BinaryEntry
- .syncSetsLocked
- IntegrationSuite
- runMonitor
- netModel
- wipe
- runNetworkMonitor
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- time.Duration
- controlServer
- TestMonitorUseCaseLifecycle
- EventType
- TrustGuard
- liftSuperseded
- exe_supersede.h
- fakeGuardRepo
- string
- pinstate.go
- LibraryClosure
- guarded_ancestor_within_limit
- github.com/spf13/cobra.Command
- TempBinary
- CopyTreeWithProgress
- NewDumpHook
- bpfstats/main.go
- btrfsLayoutFrom
- tempRow
- uninstall_system.go
- planUpdaters
- process_vm_readv.c
- github.com/cilium/ebpf.Map
- bufio.Reader
- auditGroupCoverage
- walkPlaced
- sync.Mutex
- net_tester/main.go
- ResolveLibraryClosure
- changelog.go
- IntegrationSuite
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- NetGuard
- fileedit_symlink_test.go
- preload_wait.c
- guardModel
- binledger.go
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- liveSession
- watchSet
- IntegrationSuite
- Well-known bypass classes
- newTestVetter
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- vaultFS
- Build & sign release assets (reusable workflow)
- secureResources
- parseGuardEvent
- Vault
- inode_dev
- app-listener Terminal Demo Recording
- headerWidget
- .setInspectors
- launchKey
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- tempGrant
- bpfLsmListed
- WebSocket /ws client
- time.Time
- tempGrantSetup
- install.sh
- stripTempJournal
- launch_rule_at
- make test-integration (rootful Docker suite)
- IntegrationSuite
- supersede_probe.c
- model
- trace-app-libs.sh
- IntegrationSuite
- raw_block_device.c
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- RawTarget
- .add
- IntegrationSuite
- mustSubkey
- planTaintOwners
- openBunTmp
- configedit_test.go
- runLockdown
- IntegrationSuite
- lockAndDeprovision
- findTopLevelWindow
- SystemPlacedInode
- watchPattern
- vouchedDevs
- io.Reader
- .runEditorHarness
- ParseSectionWhitelist
- graphify-refresh.sh
- Execute
- ConfigureSharedPinning

## God Nodes (most connected - your core abstractions)
1. `Guard` - 93 edges
2. `Config` - 91 edges
3. `Load()` - 85 edges
4. `IntegrationSuite` - 75 edges
5. `absPath()` - 73 edges
6. `IntegrationSuite` - 66 edges
7. `writeConfig()` - 65 edges
8. `fileEditModel` - 58 edges
9. `User` - 56 edges
10. `guardUnitTest` - 52 edges

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

## Communities (165 total, 29 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (31): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+23 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (112): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+104 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.15
Nodes (13): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+5 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): updateCatalogConfig(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+70 more)

### Community 5 - "update.go"
Cohesion: 0.06
Nodes (50): applyUpdate(), compareStableVersions(), confirmUpdate(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), fetchReleases(), filterChannel() (+42 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.08
Nodes (59): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+51 more)

### Community 7 - "shellrc.go"
Cohesion: 0.09
Nodes (53): trustManager, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), catalogPatterns(), inHome(), addKeysToAgent(), rcTarget, RefreshOptions (+45 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (19): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), addRecorder(), copySelf(), runTool(), walkEntries() (+11 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "auth.go"
Cohesion: 0.06
Nodes (54): ParseEventsFlag(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+46 more)

### Community 11 - "Guard"
Cohesion: 0.07
Nodes (11): binaryVerifyState, deferredBinary, rootHandle, sharedHashEntry, SharedPinDegraded(), canonicalBinaryPath(), eventMask(), Guard (+3 more)

### Community 12 - "Monitor"
Cohesion: 0.05
Nodes (8): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), Monitor, monitorUnitTest, pathCache, pathCacheEntry

### Community 13 - "uninstall.go"
Cohesion: 0.06
Nodes (52): confirmMasterKeyOverwrite(), runGenKey(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty() (+44 more)

### Community 14 - "install.go"
Cohesion: 0.07
Nodes (59): reloadDaemonForPasswordChange(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), restoreDaemonAfterDiffAbort() (+51 more)

### Community 15 - "runDaemon"
Cohesion: 0.14
Nodes (20): backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), makeReloadHandler(), notifySystemdReady(), reloadOnce(), reloadSlotsNeeded() (+12 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "GenerateConf"
Cohesion: 0.06
Nodes (48): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestGeneratedLibraryBlockRoundTrips(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig() (+40 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "Candidate"
Cohesion: 0.12
Nodes (27): appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), cleanOrphanedFscrypt(), addManualDirectories() (+19 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (11): ownerPaths(), Resource, GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots() (+3 more)

### Community 23 - "CandidateDir"
Cohesion: 0.06
Nodes (32): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, GlobChild, compileGlobNames(), globText() (+24 more)

### Community 24 - "go_pkg_os"
Cohesion: 0.30
Nodes (3): rejectSymlinkedComponents(), validateWatchTarget(), hint

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (44): add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector(), discover_guarded_parent() (+36 more)

### Community 27 - "NetEventType"
Cohesion: 0.11
Nodes (14): NetBpfEvent, NetEvent, NetEventType, Cstr(), FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType() (+6 more)

### Community 28 - "SanitizeText"
Cohesion: 0.08
Nodes (18): newBinaryVetter(), openBinaryVetter(), logRefreshChange(), TestLogRefreshChange(), candidates(), confirmable(), confirmAll(), describe() (+10 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.11
Nodes (27): buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), selfGuards (+19 more)

### Community 30 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "engine"
Cohesion: 0.14
Nodes (6): TestNearestRootClaims(), TestTieRank(), engine, nearestRootClaims(), pathWithin(), tieRank()

### Community 33 - "os.File"
Cohesion: 0.10
Nodes (18): refusedReplacements, engine, Guard, supersedeUnnamed(), genericElectron(), runtimeClass(), scanRuntimeClass(), TestRuntimeClassGenericElectron() (+10 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.14
Nodes (28): OpenRegularNoFollow(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed() (+20 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.08
Nodes (31): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+23 more)

### Community 36 - "Config"
Cohesion: 0.09
Nodes (42): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager, startTrustGuard(), trustStartupError(), inspectorAdmitter, watchGroup, addResource() (+34 more)

### Community 37 - "runNetworkGuard"
Cohesion: 0.09
Nodes (23): CheckBPFLSM(), CheckEBPF(), newPinGeneration(), prepareDaemonStart(), runBPFCheck(), runOneShotMode(), computeGuardBinaries(), confirmUnsafe() (+15 more)

### Community 38 - "GuardEvent"
Cohesion: 0.08
Nodes (15): holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), commMatchesGuardedBinary(), GuardEvent (+7 more)

### Community 39 - "KernelDev"
Cohesion: 0.09
Nodes (32): fdStat, BinaryStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), openRoot(), physicalParent() (+24 more)

### Community 40 - "backups.go"
Cohesion: 0.14
Nodes (23): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), Delete(), Find(), find(), Backup (+15 more)

### Community 41 - "codeedit.go"
Cohesion: 0.16
Nodes (29): nextConfig(), class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+21 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "New"
Cohesion: 0.15
Nodes (23): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails() (+15 more)

### Community 44 - "main_test.go"
Cohesion: 0.06
Nodes (18): IntegrationSuite, exploitTest, IntegrationSuite, guardBinaryFlag(), infraContainerPaths(), newEventTypes(), tailLast(), TestIntegrationSuite() (+10 more)

### Community 45 - "system.go"
Cohesion: 0.08
Nodes (37): checkAndApply(), pathCovered(), TestPathCovered(), TestDaemonUnitWithMetadataOutput(), isInsidePath(), askSSHAgentUsers(), daemonUnitWithMetadataOutput(), dupFD() (+29 more)

### Community 46 - "BinaryEntry"
Cohesion: 0.12
Nodes (17): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), VettedInode, BinaryEntry, SectionScan, BinariesSummary() (+9 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.16
Nodes (11): resInfo, engine, intersectAllows(), isSubset(), planMemberRows(), setPaths(), syncSetMap(), syncU32Map() (+3 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "runMonitor"
Cohesion: 0.16
Nodes (10): runHeadless(), runMonitor(), runTUI(), MonitorRepository, NewModel(), assertViewFits(), TestViewsFitReportedTerminalSize(), MonitorUseCase (+2 more)

### Community 51 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 52 - "wipe"
Cohesion: 0.20
Nodes (20): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+12 more)

### Community 53 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (10): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+2 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (10): diffModel, editorModel, clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), netGuardEventMsg, noticeModel, servedModel (+2 more)

### Community 56 - "time.Duration"
Cohesion: 0.10
Nodes (10): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity() (+2 more)

### Community 57 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 58 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.29
Nodes (9): NewGuardUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkMonitorUseCaseLifecycle() (+1 more)

### Community 59 - "EventType"
Cohesion: 0.10
Nodes (12): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryEventsReadOnlyRejectsRestriction(), Run(), DecodeBpfEvent() (+4 more)

### Community 60 - "TrustGuard"
Cohesion: 0.11
Nodes (8): trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 61 - "liftSuperseded"
Cohesion: 0.18
Nodes (11): dropSupersededLocked(), exeGone(), Guard, liftSuperseded(), PruneSuperseded(), supersededMark(), supersedeMaps(), supersedePruneLoop() (+3 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.18
Nodes (17): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+9 more)

### Community 64 - "string"
Cohesion: 0.14
Nodes (4): denied(), main(), denied(), main()

### Community 65 - "pinstate.go"
Cohesion: 0.23
Nodes (17): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive() (+9 more)

### Community 66 - "LibraryClosure"
Cohesion: 0.13
Nodes (14): ldConfParser, defaultLibDirs(), elfInterp(), isScript(), ldConfInclude(), ldSoConfDirs(), LibraryClosure(), resolveSoname() (+6 more)

### Community 67 - "guarded_ancestor_within_limit"
Cohesion: 0.24
Nodes (19): evict_inode_from_guard(), get_dentry_from_path(), get_inode_from_path(), guard_inode_getattr(), guard_path_link(), guard_path_mkdir(), guard_path_mknod(), guard_path_rename() (+11 more)

### Community 68 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 69 - "TempBinary"
Cohesion: 0.14
Nodes (17): confirmUserPlaced(), forwardRequest, TempBinary, TemporaryGrant, TempRule, ResolveTempBinary(), applyTempGrants(), closeTempBinaries() (+9 more)

### Community 71 - "CopyTreeWithProgress"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 73 - "bpfstats/main.go"
Cohesion: 0.19
Nodes (16): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), inertProgram(), kernelHasFunc(), TestKernelHasFunc(), trustSpec() (+8 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "tempRow"
Cohesion: 0.29
Nodes (5): exeRow, tempMaskOp, tempRow, Guard, planTempBlock()

### Community 76 - "uninstall_system.go"
Cohesion: 0.18
Nodes (17): detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertBunLaunchers(), revertSSHAgents() (+9 more)

### Community 78 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (10): copy(), is_verb(), main(), usage(), dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd() (+2 more)

### Community 80 - "github.com/cilium/ebpf.Map"
Cohesion: 0.13
Nodes (9): statFunc, deleteInoKeys(), deleteResKeys(), resolveBits(), statKey(), syncMap(), deleteKeys(), TrustGuard (+1 more)

### Community 81 - "bufio.Reader"
Cohesion: 0.19
Nodes (11): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+3 more)

### Community 82 - "auditGroupCoverage"
Cohesion: 0.18
Nodes (13): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+5 more)

### Community 83 - "walkPlaced"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 84 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 85 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 86 - "ResolveLibraryClosure"
Cohesion: 0.20
Nodes (10): fileExists(), LdSoPreloadPaths(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestMountVouchesOwnershipJudgesHostNamespace(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch(), TestResolveLibraryClosureScript() (+2 more)

### Community 87 - "changelog.go"
Cohesion: 0.21
Nodes (9): changelogText(), newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes() (+1 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "NetGuard"
Cohesion: 0.16
Nodes (9): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), eventsetSummary(), eventTypeKey(), NewNetGuard() (+1 more)

### Community 92 - "fileedit_symlink_test.go"
Cohesion: 0.35
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 94 - "guardModel"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 95 - "binledger.go"
Cohesion: 0.19
Nodes (10): grantNewBinaries(), EnsurePlaceholders(), JournalPath(), Lines(), Open(), TestEnsurePlaceholdersNoInstallDir(), TestLedgerKeepsJournalFile(), TestLedgerRecordLookupPending() (+2 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 100 - "watchSet"
Cohesion: 0.21
Nodes (6): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "newTestVetter"
Cohesion: 0.42
Nodes (9): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSwapWithinGenerationLosesRights() (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "vaultFS"
Cohesion: 0.12
Nodes (17): checkFreshDir(), openDirNoFollow(), renewDir(), TestCleanupStalePins(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode() (+9 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "secureResources"
Cohesion: 0.15
Nodes (13): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestCollectFilesystemPrereqsNoPanic() (+5 more)

### Community 110 - "parseGuardEvent"
Cohesion: 0.22
Nodes (7): bpfGuardEvent, fsGateLabel(), parseGuardEvent(), processGateLabel(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels(), suspectLabel()

### Community 111 - "Vault"
Cohesion: 0.24
Nodes (6): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata()

### Community 112 - "inode_dev"
Cohesion: 0.31
Nodes (6): ctx_ptr(), inode_dev(), sb_dev(), guard_bprm_committed(), guard_inode_free(), trust_inode_free()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "headerWidget"
Cohesion: 0.18
Nodes (6): columnState, dataRow, headerWidget, newColumnState(), newDataRow(), newHeaderWidget()

### Community 115 - ".setInspectors"
Cohesion: 0.33
Nodes (4): engine, Inspector, TrustGuard, SetInspectors()

### Community 116 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "tempGrant"
Cohesion: 0.21
Nodes (4): tempGrant, loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 121 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "time.Time"
Cohesion: 0.25
Nodes (3): guiModel, Run(), exeVerdict

### Community 124 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "stripTempJournal"
Cohesion: 0.16
Nodes (14): ensurePinStateFilePlaceholder(), ensurePlaceholder(), writeInPlace(), bootstrapSelfProtectPlaceholders(), ensureHashFilePlaceholder(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace() (+6 more)

### Community 127 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "model"
Cohesion: 0.26
Nodes (6): formatType(), listenForEvents(), tickStats(), eventLine, eventMsg, model

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "raw_block_device.c"
Cohesion: 0.83
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 143 - ".add"
Cohesion: 0.23
Nodes (3): TemporaryAccess, fakeTempGrant, tempTrace

### Community 145 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 147 - "planTaintOwners"
Cohesion: 0.27
Nodes (8): taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), planTaintOwners(), setKeyOf(), allowSet(), TestPlanTaintOwners(), TestPlanTaintOwners_ReloadOverlapKeepsSetKey()

### Community 148 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 149 - "configedit_test.go"
Cohesion: 0.57
Nodes (7): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession()

### Community 150 - "runLockdown"
Cohesion: 0.29
Nodes (8): configureDaemonLogging(), resolveConfigPath(), runLockdown(), dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 152 - "lockAndDeprovision"
Cohesion: 0.33
Nodes (6): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 153 - "findTopLevelWindow"
Cohesion: 0.38
Nodes (4): findTopLevelWindow(), floatWindow(), internAtom(), setFloatHint()

### Community 154 - "SystemPlacedInode"
Cohesion: 0.29
Nodes (7): hostView(), mountVouchesOwnership(), rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), TestHostView(), SystemPlacedInode()

### Community 155 - "watchPattern"
Cohesion: 0.60
Nodes (5): addSystemBinary(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern

### Community 156 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 157 - "io.Reader"
Cohesion: 0.40
Nodes (4): entryOf(), BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype()

### Community 159 - "ParseSectionWhitelist"
Cohesion: 0.50
Nodes (4): isLibDirectiveText(), IsLibraryDirective(), ParseSectionWhitelist(), TestParseSectionWhitelist()

## Knowledge Gaps
- **52 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+47 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 365 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **29 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.061) - this node is a cross-community bridge._
- **Why does `guardDeltaEvents()` connect `IntegrationSuite` to `absPath`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Why does `fileEditModel` connect `fileEditModel` to `fileedit_model_test.go`, `codeedit.go`, `vaultFS`, `github.com/charmbracelet/bubbletea.Cmd`, `go_pkg_os`, `fileedit_symlink_test.go`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _52 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11625277766948656 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.026083112290008842 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._