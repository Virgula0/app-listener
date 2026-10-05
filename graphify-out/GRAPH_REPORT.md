# Graph Report - app-listener  (2026-10-05)

## Corpus Check
- 314 files · ~1,487,895 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4160 nodes · 16523 edges · 172 communities (142 shown, 30 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1214 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d633567d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- wipe
- daemonconfig_test.go
- update.go
- monitor.bpf.c
- User
- guardUnitTest
- globreserve_test.go
- auth_test.go
- Guard
- monitorUnitTest
- runUninstall
- applyDiffAdditions
- runDaemon
- IntegrationSuite
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- trustManager
- usecase/daemon_test.go
- Resource
- IntegrationSuite
- Candidate
- guard.bpf.c
- dataRowRenderer
- NetEventType
- Ledger
- buildOneGuard
- Monitor
- serve.go
- engine
- os.File
- fileedit_model_test.go
- DaemonEvent
- Config
- runNetworkMonitor
- GuardEvent
- NetworkMonitor
- codeedit.go
- highlight_test.go
- absPath
- SelectOrphans
- IntegrationSuite
- safeio.go
- netModel
- .syncSetsLocked
- IntegrationSuite
- bpfstats/main.go
- runForward
- go_pkg_os
- SanitizeText
- github.com/charmbracelet/bubbletea.Cmd
- catalogRefresher
- controlServer
- TestMonitorUseCaseLifecycle
- EventType
- TrustGuard
- liftSuperseded
- exe_supersede.h
- fakeGuardRepo
- string
- pinstate.go
- watchSet
- guarded_ancestor_within_limit
- NewDumpHook
- headerWidget
- copytree.go
- bufio.Reader
- RefreshSection
- btrfsLayoutFrom
- time.Time
- securewrite.go
- planUpdaters
- process_vm_readv.c
- syncMap
- install.go
- IntegrationSuite
- .step
- sync.Mutex
- net_tester/main.go
- render.go
- audit.go
- IntegrationSuite
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- go_pkg_path_filepath
- runGuard
- preload_wait.c
- procstats.go
- time.Duration
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- liveSession
- IntegrationSuite
- IntegrationSuite
- Well-known bypass classes
- newTestVetter
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- vaultFS
- Build & sign release assets (reusable workflow)
- Vault
- RawTarget
- Vault
- inode_dev
- app-listener Terminal Demo Recording
- expandPlaceholders
- .setInspectors
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- github.com/cilium/ebpf.Map
- bpfLsmListed
- WebSocket /ws client
- github.com/spf13/cobra.Command
- watchPattern
- install.sh
- IntegrationSuite
- launch_rule_at
- make test-integration (rootful Docker suite)
- lockAndDeprovision
- supersede_probe.c
- CandidateDir
- trace-app-libs.sh
- IntegrationSuite
- raw_block_device.c
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- Open
- BinaryEntry
- runNetworkGuard
- New
- SuperblockDev
- TempBinary
- bottomBar
- configedit_test.go
- launchKey
- .runEditorHarness
- .write
- .resyncLink
- IntegrationSuite
- find_test.go
- io.Writer
- findTopLevelWindow
- fakeDaemon
- Vault
- graphify-refresh.sh
- findSectionStart
- ParseEventsFlag
- .pathInfo
- resolveBits
- guardLSMHooks
- github.com/google/fscrypt/actions.Context
- formatEntry
- harnessSuite
- TestMain
- parsePasswd
- pathCovered

## God Nodes (most connected - your core abstractions)
1. `Guard` - 93 edges
2. `Config` - 91 edges
3. `Load()` - 85 edges
4. `IntegrationSuite` - 75 edges
5. `absPath()` - 73 edges
6. `IntegrationSuite` - 66 edges
7. `writeConfig()` - 65 edges
8. `fileEditModel` - 60 edges
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

## Communities (172 total, 30 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.08
Nodes (98): AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), confirmUserPlaced(), init(), init(), init(), init() (+90 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (116): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+108 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "wipe"
Cohesion: 0.20
Nodes (10): isFileVaultCiphertext(), isRegularFileTarget(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall(), TestNewBoundedKeyFnRetryAborts() (+2 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (75): updateCatalogConfig(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+67 more)

### Community 5 - "update.go"
Cohesion: 0.06
Nodes (46): changelogText(), showChangelog(), TestChangelogText(), compareStableVersions(), confirmUpdate(), downloadAndVerify(), downloadFile(), downloadReleaseFiles() (+38 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.08
Nodes (59): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+51 more)

### Community 7 - "User"
Cohesion: 0.12
Nodes (39): trustManager, bunEntryConfigured(), bunLaunchersFor(), addKeysToAgent(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent() (+31 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), addRecorder(), copySelf(), runTool(), walkEntries() (+13 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.17
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "auth_test.go"
Cohesion: 0.13
Nodes (29): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+21 more)

### Community 11 - "Guard"
Cohesion: 0.08
Nodes (10): binaryVerifyState, rootHandle, sharedHashEntry, openRoot(), SharedPinDegraded(), eventMask(), Guard, hashBinaryShared() (+2 more)

### Community 12 - "monitorUnitTest"
Cohesion: 0.10
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 13 - "runUninstall"
Cohesion: 0.09
Nodes (30): deletePostBackups(), restoreBackups(), askDecrypt(), decryptDirectories(), decryptStep(), offerBackupCleanup(), pickDirsToDecrypt(), runUninstall() (+22 more)

### Community 14 - "applyDiffAdditions"
Cohesion: 0.14
Nodes (19): askBunTmpdirUsers(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), configPaths(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), promptEditPassword() (+11 more)

### Community 15 - "runDaemon"
Cohesion: 0.07
Nodes (40): CheckBPFLSM(), CheckEBPF(), awaitStartupOrSignal(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), makeReloadHandler() (+32 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "GenerateConf"
Cohesion: 0.07
Nodes (45): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestGeneratedLibraryBlockRoundTrips(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig() (+37 more)

### Community 19 - "fileEditModel"
Cohesion: 0.10
Nodes (9): clipLabel(), humanSize(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "trustManager"
Cohesion: 0.11
Nodes (22): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), trustManager (+14 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.11
Nodes (54): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+46 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (11): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), uniqueEncryptionRoots() (+3 more)

### Community 23 - "IntegrationSuite"
Cohesion: 0.23
Nodes (7): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 24 - "Candidate"
Cohesion: 0.09
Nodes (35): setConfinedHomes(), appendSectionsAndEdit(), collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), cleanOrphanedFscrypt(), addManualDirectories() (+27 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (44): add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector(), discover_guarded_parent() (+36 more)

### Community 27 - "NetEventType"
Cohesion: 0.12
Nodes (12): NetBpfEvent, NetEventType, Cstr(), FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), eventsetSummary() (+4 more)

### Community 28 - "Ledger"
Cohesion: 0.15
Nodes (3): Ledger, Pending, Source

### Community 29 - "buildOneGuard"
Cohesion: 0.11
Nodes (27): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan() (+19 more)

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "engine"
Cohesion: 0.11
Nodes (6): TestNearestRootClaims(), TestTieRank(), engine, nearestRootClaims(), pathWithin(), tieRank()

### Community 33 - "os.File"
Cohesion: 0.08
Nodes (29): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+21 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.09
Nodes (44): AtomicWriteAt(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth() (+36 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.13
Nodes (25): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+17 more)

### Community 36 - "Config"
Cohesion: 0.11
Nodes (40): ServeConfig, newDaemonModel(), runDaemonUI(), runServedTUI(), runTUI(), watchGroup, addResource(), applyAllowLib() (+32 more)

### Community 37 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (10): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+2 more)

### Community 38 - "GuardEvent"
Cohesion: 0.07
Nodes (18): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, fsGateLabel(), GuardEvent, logBacklogDenial(), parseGuardEvent() (+10 more)

### Community 39 - "NetworkMonitor"
Cohesion: 0.11
Nodes (9): NetEvent, foreignPinnedLink(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), newFakeNetworkMonitorRepo(), NetworkMonitor (+1 more)

### Community 40 - "codeedit.go"
Cohesion: 0.25
Nodes (16): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+8 more)

### Community 41 - "highlight_test.go"
Cohesion: 0.10
Nodes (32): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep() (+24 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "SelectOrphans"
Cohesion: 0.30
Nodes (13): appProtectorSet(), isAppProtector(), orphanedPolicies(), orphanedProtectors(), SelectOrphans(), appProtector(), foreignProtector(), policy() (+5 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.09
Nodes (7): IntegrationSuite, exploitTest, IntegrationSuite, newEventTypes(), tailLast(), IntegrationSuite, editorHarness

### Community 45 - "safeio.go"
Cohesion: 0.07
Nodes (38): checkAndApply(), TestDaemonUnitWithMetadataOutput(), daemonUnitWithMetadataOutput(), dupFD(), installFile(), installFileAs(), installServices(), installSSHAgent() (+30 more)

### Community 46 - "netModel"
Cohesion: 0.10
Nodes (16): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), formatAddr() (+8 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.11
Nodes (18): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, intersectAllows(), isSubset(), planMemberRows() (+10 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "bpfstats/main.go"
Cohesion: 0.19
Nodes (16): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), inertProgram(), kernelHasFunc(), TestKernelHasFunc(), trustSpec() (+8 more)

### Community 51 - "runForward"
Cohesion: 0.12
Nodes (18): dialLiveSession(), applyEdit(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), chooseResources() (+10 more)

### Community 52 - "go_pkg_os"
Cohesion: 0.08
Nodes (18): Execute(), rejectSymlinkedComponents(), validateWatchTarget(), classifySetupError(), isNotEncryptedErr(), readKeyFrom(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup() (+10 more)

### Community 53 - "SanitizeText"
Cohesion: 0.29
Nodes (4): ownerPaths(), binaryVetter, vetDecision, SanitizeText()

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (19): diffModel, editorModel, listenForDaemonEvents(), NewDaemonModel(), tickDaemonStats(), syncViewport(), renderBullets(), TestRenderBulletsVerbatimAndIndented() (+11 more)

### Community 56 - "catalogRefresher"
Cohesion: 0.18
Nodes (6): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), catalogRefresher

### Community 57 - "controlServer"
Cohesion: 0.19
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 58 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.27
Nodes (10): NewGuardUseCase(), NewMonitorUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError(), TestNetworkGuardUseCaseLifecycle() (+2 more)

### Community 59 - "EventType"
Cohesion: 0.07
Nodes (17): systemPatterns(), systemRule(), runHeadless(), BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents() (+9 more)

### Community 60 - "TrustGuard"
Cohesion: 0.14
Nodes (5): trustRows, setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 61 - "liftSuperseded"
Cohesion: 0.15
Nodes (10): dropSupersededLocked(), exeGone(), Guard, SupersededBinary, liftSuperseded(), PruneSuperseded(), supersededMark(), supersedeMaps() (+2 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.18
Nodes (17): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+9 more)

### Community 64 - "string"
Cohesion: 0.14
Nodes (4): denied(), main(), denied(), main()

### Community 65 - "pinstate.go"
Cohesion: 0.12
Nodes (29): configureDaemonLogging(), loadDaemonConfig(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState() (+21 more)

### Community 66 - "watchSet"
Cohesion: 0.24
Nodes (6): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet

### Community 67 - "guarded_ancestor_within_limit"
Cohesion: 0.24
Nodes (19): evict_inode_from_guard(), get_dentry_from_path(), get_inode_from_path(), guard_inode_getattr(), guard_path_link(), guard_path_mkdir(), guard_path_mknod(), guard_path_rename() (+11 more)

### Community 68 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 69 - "headerWidget"
Cohesion: 0.19
Nodes (5): columnState, dataRow, headerWidget, newDataRow(), newHeaderWidget()

### Community 71 - "copytree.go"
Cohesion: 0.20
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 73 - "RefreshSection"
Cohesion: 0.17
Nodes (14): logRefreshChange(), TestLogRefreshChange(), RefreshOptions, SectionScan, catalogEntryMatch(), findCatalogRoot(), findCatalogWatchSubPath(), isInsidePath() (+6 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "time.Time"
Cohesion: 0.24
Nodes (4): guiModel, newColumnState(), Run(), exeVerdict

### Community 76 - "securewrite.go"
Cohesion: 0.31
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 78 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (10): copy(), is_verb(), main(), usage(), dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd() (+2 more)

### Community 81 - "install.go"
Cohesion: 0.05
Nodes (54): confirmMasterKeyOverwrite(), runGenKey(), reloadDaemonForPasswordChange(), applyLiveRefresh(), buildBinaryIfNeeded(), ensureInstalledBinary(), ensureMasterKey(), installBinaryOnly() (+46 more)

### Community 83 - ".step"
Cohesion: 0.36
Nodes (5): placedWalk, readlinkFd(), rootOwnedStat(), systemFileStat(), walkPlaced()

### Community 84 - "sync.Mutex"
Cohesion: 0.21
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 85 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 86 - "render.go"
Cohesion: 0.15
Nodes (13): lineStyles, viewportState, baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 87 - "audit.go"
Cohesion: 0.23
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "go_pkg_path_filepath"
Cohesion: 0.37
Nodes (3): LiveEmptyWhitelistRejected(), TestLiveEmptyWhitelistRejected(), hint

### Community 92 - "runGuard"
Cohesion: 0.22
Nodes (9): selfProtectSpecs(), resolveGuardConfig(), runGuard(), runGuardTUI(), selfProtectSpec, Mode, guardModeKey(), modeString() (+1 more)

### Community 94 - "procstats.go"
Cohesion: 0.16
Nodes (14): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+6 more)

### Community 95 - "time.Duration"
Cohesion: 0.21
Nodes (4): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

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
Cohesion: 0.31
Nodes (4): applyNewFileMeta(), listEntries(), dirEntry, vaultFS

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "Vault"
Cohesion: 0.15
Nodes (17): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq() (+9 more)

### Community 110 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 111 - "Vault"
Cohesion: 0.23
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 112 - "inode_dev"
Cohesion: 0.31
Nodes (6): ctx_ptr(), inode_dev(), sb_dev(), guard_bprm_committed(), guard_inode_free(), trust_inode_free()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "expandPlaceholders"
Cohesion: 0.10
Nodes (16): symlinkTargetGlob(), confSafeMatch(), expandPlaceholders(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), CandidateDir, ConfSafePath() (+8 more)

### Community 115 - ".setInspectors"
Cohesion: 0.46
Nodes (3): engine, Inspector, SetInspectors()

### Community 116 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.15
Nodes (11): findButton, span, column(), moveTo(), findAll(), Finder, lowerRunes(), newFindInput() (+3 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "github.com/cilium/ebpf.Map"
Cohesion: 0.21
Nodes (7): deleteInoKeys(), deleteResKeys(), syncU32Map(), loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows(), deleteKeys()

### Community 121 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "github.com/spf13/cobra.Command"
Cohesion: 0.27
Nodes (11): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink(), effectiveVerbose() (+3 more)

### Community 124 - "watchPattern"
Cohesion: 0.33
Nodes (7): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 127 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "lockAndDeprovision"
Cohesion: 0.27
Nodes (7): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "CandidateDir"
Cohesion: 0.12
Nodes (16): entryWriters(), reserveChain(), underResource(), globBuilder, GlobChild, compileGlobNames(), globText(), GlobReservations (+8 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "raw_block_device.c"
Cohesion: 0.83
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "Open"
Cohesion: 0.25
Nodes (8): newBinaryVetter(), openBinaryVetter(), EnsurePlaceholders(), JournalPath(), Open(), TestEnsurePlaceholdersNoInstallDir(), TestLedgerKeepsJournalFile(), TestLedgerRecordLookupPending()

### Community 143 - "BinaryEntry"
Cohesion: 0.11
Nodes (15): deferredBinary, VettedInode, BinaryEntry, BinariesSummary(), canonicalBinaryPath(), canonicalPaths(), commMatchesGuardedBinary(), ConfinedEntry() (+7 more)

### Community 144 - "runNetworkGuard"
Cohesion: 0.11
Nodes (17): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode, NetGuardEvent (+9 more)

### Community 145 - "New"
Cohesion: 0.07
Nodes (54): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath() (+46 more)

### Community 146 - "SuperblockDev"
Cohesion: 0.17
Nodes (10): checkBtrfsKeys(), entryOf(), FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev(), SuperblockDevPath(), TestMountinfoDev() (+2 more)

### Community 147 - "TempBinary"
Cohesion: 0.09
Nodes (21): exeRow, tempGrant, tempMaskOp, tempRow, Guard, TempBinary, TemporaryGrant, TempRule (+13 more)

### Community 148 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 149 - "configedit_test.go"
Cohesion: 0.57
Nodes (7): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession()

### Community 150 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 152 - ".write"
Cohesion: 0.22
Nodes (8): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), writeFileInPlace()

### Community 153 - ".resyncLink"
Cohesion: 0.39
Nodes (3): engine, Guard, supersedeUnnamed()

### Community 155 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 156 - "io.Writer"
Cohesion: 0.52
Nodes (6): candidates(), confirmable(), confirmAll(), describe(), run(), candidate

### Community 157 - "findTopLevelWindow"
Cohesion: 0.38
Nodes (4): findTopLevelWindow(), floatWindow(), internAtom(), setFloatHint()

### Community 161 - "findSectionStart"
Cohesion: 0.29
Nodes (7): findSectionStart(), isLibDirectiveText(), ParseSectionHeaderPath(), unquotePath(), IsLibraryDirective(), ParseSectionWhitelist(), TestParseSectionWhitelist()

### Community 162 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 163 - ".pathInfo"
Cohesion: 0.48
Nodes (3): evalSymlinksOrEmpty(), pathCache, pathCacheEntry

### Community 164 - "resolveBits"
Cohesion: 0.43
Nodes (4): statFunc, resolveBits(), statKey(), UpdaterPlan

### Community 165 - "guardLSMHooks"
Cohesion: 0.33
Nodes (3): trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent()

### Community 166 - "github.com/google/fscrypt/actions.Context"
Cohesion: 0.40
Nodes (4): modifiedContextWithSource(), CleanOrphans(), listPolicies(), listProtectors()

### Community 167 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 169 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 170 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 171 - "pathCovered"
Cohesion: 0.67
Nodes (3): pathCovered(), TestPathCovered(), isInsidePath()

## Knowledge Gaps
- **52 isolated node(s):** `hint`, `tempJournalRow`, `launchRule`, `trustEvent`, `btrfsIoctlFsInfoArgs` (+47 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 375 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **30 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `go_pkg_strings`, `IntegrationSuite`?**
  _High betweenness centrality (0.052) - this node is a cross-community bridge._
- **Why does `shQuote()` connect `IntegrationSuite` to `daemon/daemon.go`, `IntegrationSuite`, `absPath`, `IntegrationSuite`, `IntegrationSuite`, `.runEditorHarness`, `IntegrationSuite`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `os.File`, `DaemonEvent`, `GuardEvent`, `time.Time`, `runDaemon`, `BinaryEntry`, `sync.Mutex`, `runGuard`, `buildOneGuard`, `fakeDaemon`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **What connects `hint`, `tempJournalRow`, `launchRule` to the rest of the system?**
  _52 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07847714824459011 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.025352513886910698 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._