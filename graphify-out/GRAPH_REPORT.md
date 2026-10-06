# Graph Report - app-listener  (2026-10-06)

## Corpus Check
- 319 files · ~1,495,294 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4212 nodes · 16940 edges · 169 communities (137 shown, 32 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1236 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `667f2414`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- wipe
- daemonconfig_test.go
- update_test.go
- monitor.bpf.c
- User
- guardUnitTest
- mkdirs
- auth.go
- testcontainers.Container
- Config
- ListUsers
- highlight_test.go
- runDaemon
- .startGuardStd
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- buildTrustedSet
- usecase/daemon_test.go
- Resource
- KernelDev
- Candidate
- guard.bpf.c
- Run
- NetEventType
- Ledger
- BinaryEntry
- SelectOrphans
- serve.go
- GuardInodeKey
- monitorUnitTest
- fileedit_model_test.go
- DaemonEvent
- watchGroup
- update.go
- GuardEvent
- SanitizeText
- RefreshSection
- TemporaryGrant
- IntegrationSuite
- check_and_emit_args
- IntegrationSuite
- os.FileMode
- Monitor
- github.com/cilium/ebpf.Map
- IntegrationSuite
- filevault.go
- IntegrationSuite
- IntegrationSuite
- Vault
- github.com/charmbracelet/bubbletea.Cmd
- absPath
- os.File
- controlServer
- updateCatalogConfig
- EventType
- TrustGuard
- .resyncLink
- exe_supersede.h
- fakeGuardRepo
- string
- pinstate.go
- LibraryClosure
- time.Duration
- github.com/spf13/cobra.Command
- catalogRefresher
- copyRegularAt
- bufio.Reader
- auditGroupCoverage
- syncMap
- filevault_test.go
- serveWebSocket
- IntegrationSuite
- planUpdaters
- process_vm_readv.c
- go_pkg_github_com_testcontainers_testcontainers_go
- runMaintenanceMode
- IntegrationSuite
- .step
- sync.Mutex
- net_tester/main.go
- render.go
- time.Time
- IntegrationSuite
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- NewEventFanout
- guardModel
- preload_wait.c
- classifySupportError
- runNetworkGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- liveSession
- NetworkMonitorUseCase
- IntegrationSuite
- Well-known bypass classes
- newTestVetter
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- runMonitor
- Build & sign release assets (reusable workflow)
- runGenKey
- daemon_multicall_test.go
- Vault
- inode_dev
- app-listener Terminal Demo Recording
- CandidateDir
- vaultFS
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- launchKey
- bpfLsmListed
- WebSocket /ws client
- NewDumpHook
- writeWithin
- install.sh
- TempBinary
- launch_rule_at
- make test-integration (rootful Docker suite)
- Guard
- supersede_probe.c
- binaryVerifyState
- trace-app-libs.sh
- IntegrationSuite
- raw_block_device.c
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- shQuote
- bottomBar
- verifyReleaseWithKey
- io.Reader
- Guard
- tempRow
- fakeDaemon
- configedit_test.go
- DiscoverInfraBinaries
- bpfstats/main.go
- lockRootRecovering
- newTestSession
- harnessSuite
- find.go
- prepareDaemonStart
- VerifyInstalledResourcesLocked
- MulticallError
- Vault
- graphify-refresh.sh
- IntegrationSuite
- StripPinnedTempAllows
- copyXattrs
- golang.org/x/sys/unix.Timespec
- main_test.go

## God Nodes (most connected - your core abstractions)
1. `Guard` - 93 edges
2. `Config` - 91 edges
3. `Load()` - 85 edges
4. `GuardInodeKey` - 82 edges
5. `IntegrationSuite` - 75 edges
6. `absPath()` - 74 edges
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

## Communities (169 total, 32 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (16): Execute(), tempJournalRow, launchRule, trustEvent, btrfsIoctlFsInfoArgs, Check(), collectDiagnostics(), readSysctl() (+8 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (123): plantBun(), TestOpenBunTmp_KeepsVettedDir(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait() (+115 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.12
Nodes (6): IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), daemonEvent

### Community 3 - "wipe"
Cohesion: 0.14
Nodes (17): deprovisionKind, isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), isLockedRegularFileErr(), isNotEncryptedErr(), newBoundedKeyFn(), readKey() (+9 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (69): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+61 more)

### Community 5 - "update_test.go"
Cohesion: 0.09
Nodes (14): downloadFile(), downloadReleaseFiles(), fetchReleases(), rsaPublicKeyPEM(), signedChecksum(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases() (+6 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.08
Nodes (59): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+51 more)

### Community 7 - "User"
Cohesion: 0.07
Nodes (55): trustManager, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), catalogPatterns(), inHome(), dupFD(), installServices(), installSSHAgent() (+47 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (23): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+15 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "auth.go"
Cohesion: 0.07
Nodes (42): ParseEventsFlag(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+34 more)

### Community 11 - "testcontainers.Container"
Cohesion: 0.17
Nodes (3): configBinaryLines(), rawExec(), watchPathsInConfig()

### Community 12 - "Config"
Cohesion: 0.08
Nodes (39): ServeConfig, newDaemonModel(), runDaemonUI(), runServedTUI(), runTUI(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager (+31 more)

### Community 13 - "ListUsers"
Cohesion: 0.08
Nodes (38): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), deletePostBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata(), decryptStep() (+30 more)

### Community 14 - "highlight_test.go"
Cohesion: 0.11
Nodes (27): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+19 more)

### Community 15 - "runDaemon"
Cohesion: 0.09
Nodes (29): awaitStartupOrSignal(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), makeReloadHandler(), notifySystemdReady(), reloadOnce() (+21 more)

### Community 16 - ".startGuardStd"
Cohesion: 0.12
Nodes (7): attrOpCase, eventsForPath(), guardDeltaEvents(), guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "GenerateConf"
Cohesion: 0.06
Nodes (54): TestDiffMergePreservesExistingAndParses(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips(), TestGroupedSectionAddressedByEncryptionRoot() (+46 more)

### Community 19 - "fileEditModel"
Cohesion: 0.10
Nodes (9): clipLabel(), humanSize(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "buildTrustedSet"
Cohesion: 0.16
Nodes (16): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.11
Nodes (54): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+46 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (11): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), uniqueEncryptionRoots() (+3 more)

### Community 23 - "KernelDev"
Cohesion: 0.08
Nodes (37): fdStat, BinaryStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD() (+29 more)

### Community 24 - "Candidate"
Cohesion: 0.12
Nodes (26): appendSectionsAndEdit(), collectDiffAdditions(), pathCovered(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories() (+18 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (52): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class() (+44 more)

### Community 26 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 27 - "NetEventType"
Cohesion: 0.07
Nodes (19): NetBpfEvent, NetEvent, NetEventType, foreignPinnedLink(), FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType() (+11 more)

### Community 28 - "Ledger"
Cohesion: 0.10
Nodes (16): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), EnsurePlaceholders(), Ledger (+8 more)

### Community 29 - "BinaryEntry"
Cohesion: 0.07
Nodes (40): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan() (+32 more)

### Community 30 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 31 - "serve.go"
Cohesion: 0.13
Nodes (22): basicAuth(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop(), relayServeSignals() (+14 more)

### Community 32 - "GuardInodeKey"
Cohesion: 0.05
Nodes (24): refusedReplacements, TestNearestRootClaims(), TestTieRank(), engine, nearestRootClaims(), pathWithin(), tieRank(), deleteInoKeys() (+16 more)

### Community 33 - "monitorUnitTest"
Cohesion: 0.10
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.08
Nodes (45): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel() (+37 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.09
Nodes (32): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+24 more)

### Community 36 - "watchGroup"
Cohesion: 0.10
Nodes (36): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+28 more)

### Community 37 - "update.go"
Cohesion: 0.15
Nodes (25): showChangelog(), compareStableVersions(), confirmUpdate(), downloadAndVerify(), filterChannel(), isTerminal(), latestTime(), newerThanInstalled() (+17 more)

### Community 38 - "GuardEvent"
Cohesion: 0.06
Nodes (27): selfProtectSpecs(), displayEntries(), holdForward(), holdHeadless(), pingInterval(), chooseResource(), runLiveEdit(), HeadlessLine() (+19 more)

### Community 39 - "SanitizeText"
Cohesion: 0.19
Nodes (7): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), binaryVetter, vetDecision, SanitizeText()

### Community 40 - "RefreshSection"
Cohesion: 0.15
Nodes (16): logRefreshChange(), TestLogRefreshChange(), RefreshOptions, SectionScan, admitNew(), SectionChange, keepExisting(), LiveEmptyWhitelistRejected() (+8 more)

### Community 41 - "TemporaryGrant"
Cohesion: 0.18
Nodes (7): TemporaryGrant, TempRule, applyTempGrants(), TemporaryAccess, planTempGrants(), revokeTempGrants(), tempJournalRows()

### Community 43 - "check_and_emit_args"
Cohesion: 0.18
Nodes (14): exe_is_superseded(), exe_refused(), check_and_emit_args(), code_suspect(), current_is_inspector(), exe_tag_bits(), fill_path(), get_current_exe_inode() (+6 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): guardBinaryFlag(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "os.FileMode"
Cohesion: 0.11
Nodes (23): TestDaemonUnitWithMetadataOutput(), daemonUnitWithMetadataOutput(), installFile(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt() (+15 more)

### Community 46 - "Monitor"
Cohesion: 0.12
Nodes (5): evalSymlinksOrEmpty(), NewMonitor(), Monitor, pathCache, pathCacheEntry

### Community 48 - "github.com/cilium/ebpf.Map"
Cohesion: 0.11
Nodes (19): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), deleteResKeys(), engine, isSubset(), planMemberRows() (+11 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "filevault.go"
Cohesion: 0.14
Nodes (21): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), Vault, newFileVaultAEAD(), openFileVault(), openFileVaultWithMasterKey(), recoverFileInPlace() (+13 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 53 - "Vault"
Cohesion: 0.17
Nodes (13): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), TestCollectFilesystemPrereqsNoPanic(), isFileVaultCiphertext() (+5 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (22): runTUI(), diffModel, editorModel, formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), NewNetModel() (+14 more)

### Community 55 - "absPath"
Cohesion: 0.12
Nodes (3): IntegrationSuite, IntegrationSuite, absPath()

### Community 56 - "os.File"
Cohesion: 0.05
Nodes (54): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), setConfinedHomes() (+46 more)

### Community 57 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 58 - "updateCatalogConfig"
Cohesion: 0.17
Nodes (15): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+7 more)

### Community 59 - "EventType"
Cohesion: 0.09
Nodes (13): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryEventsReadOnlyRejectsRestriction(), planTempBlock(), Run() (+5 more)

### Community 60 - "TrustGuard"
Cohesion: 0.13
Nodes (6): trustHook, trustRows, setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 61 - ".resyncLink"
Cohesion: 0.12
Nodes (13): engine, Guard, supersedeUnnamed(), dropSupersededLocked(), exeGone(), Guard, liftSuperseded(), PruneSuperseded() (+5 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.20
Nodes (16): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+8 more)

### Community 64 - "string"
Cohesion: 0.14
Nodes (4): denied(), main(), denied(), main()

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (23): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+15 more)

### Community 66 - "LibraryClosure"
Cohesion: 0.12
Nodes (16): ldConfParser, defaultLibDirs(), elfInterp(), isScript(), ldConfInclude(), ldSoConfDirs(), LibraryClosure(), rootOwnedInode() (+8 more)

### Community 68 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 69 - "catalogRefresher"
Cohesion: 0.14
Nodes (10): drain(), newCatalogRefresher(), readEvents(), sameConfig(), addSystemBinary(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+2 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 72 - "bufio.Reader"
Cohesion: 0.15
Nodes (13): swapFile(), controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines() (+5 more)

### Community 73 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 74 - "syncMap"
Cohesion: 0.06
Nodes (24): GlobChild, statFunc, BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), compileGlobNames(), globText() (+16 more)

### Community 75 - "filevault_test.go"
Cohesion: 0.26
Nodes (21): EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), stageRecovery(), inode(), TestEnsureRecoverySidecarPlaceholder(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile() (+13 more)

### Community 76 - "serveWebSocket"
Cohesion: 0.21
Nodes (7): configureWebSocket(), sameOrigin(), serveWebSocket(), TestSameOrigin(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 78 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (10): copy(), is_verb(), main(), usage(), dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd() (+2 more)

### Community 80 - "go_pkg_github_com_testcontainers_testcontainers_go"
Cohesion: 0.19
Nodes (3): inNosuidNamespace(), launchCase, pooledContainer

### Community 81 - "runMaintenanceMode"
Cohesion: 0.06
Nodes (45): reloadDaemonForPasswordChange(), applyLiveRefresh(), buildBinaryIfNeeded(), installBinaryOnly(), parseMaintenanceFlags(), preflightDeployedBinary(), runMaintenanceMode(), runUpdateCatalogOnly() (+37 more)

### Community 83 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 84 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 85 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 86 - "render.go"
Cohesion: 0.15
Nodes (13): lineStyles, viewportState, baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 87 - "time.Time"
Cohesion: 0.23
Nodes (7): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, dirEntry

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "NewEventFanout"
Cohesion: 0.22
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 92 - "guardModel"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 94 - "classifySupportError"
Cohesion: 0.18
Nodes (10): classifySetupError(), classifySupportError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), TestClassifySupportErrorEncryptionNotEnabled(), TestClassifySupportErrorEncryptionNotEnabledF2fs(), TestClassifySupportErrorGeneric() (+2 more)

### Community 95 - "runNetworkGuard"
Cohesion: 0.11
Nodes (19): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), runHeadless() (+11 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 100 - "NetworkMonitorUseCase"
Cohesion: 0.12
Nodes (13): NetworkMonitorRepository, NewGuardUseCase(), NetworkMonitorUseCase, NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle() (+5 more)

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

### Community 107 - "runMonitor"
Cohesion: 0.11
Nodes (19): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+11 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runGenKey"
Cohesion: 0.22
Nodes (10): confirmMasterKeyOverwrite(), runGenKey(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen(), GenerateMasterKey() (+2 more)

### Community 110 - "daemon_multicall_test.go"
Cohesion: 0.40
Nodes (5): decodeMulticall(), IntegrationSuite, multicallConfig(), uutilsLinks(), mcProbe

### Community 111 - "Vault"
Cohesion: 0.16
Nodes (12): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), CopyTreeWithProgress(), TestCopyTreeProgress() (+4 more)

### Community 112 - "inode_dev"
Cohesion: 0.27
Nodes (7): ctx_ptr(), inode_dev(), sb_dev(), guard_exec_applet(), guard_inode_free(), mc_basename(), trust_inode_free()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "CandidateDir"
Cohesion: 0.07
Nodes (26): systemPatterns(), systemRule(), entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName() (+18 more)

### Community 115 - "vaultFS"
Cohesion: 0.36
Nodes (3): applyNewFileMeta(), listEntries(), vaultFS

### Community 116 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.14
Nodes (20): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+12 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 121 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 124 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "TempBinary"
Cohesion: 0.26
Nodes (10): TempBinary, ResolveTempBinary(), TestResolveTempBinaryClassesRuntime(), closeTempBinaries(), daemonUseCase, logTempGrant(), newExeTap(), resolveTempBinaries() (+2 more)

### Community 127 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "binaryVerifyState"
Cohesion: 0.43
Nodes (3): binaryVerifyState, sharedHashEntry, hashBinaryShared()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "raw_block_device.c"
Cohesion: 0.83
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "shQuote"
Cohesion: 0.22
Nodes (3): IntegrationSuite, shQuote(), IntegrationSuite

### Community 143 - "bottomBar"
Cohesion: 0.23
Nodes (7): decryptDirectories(), TestDecryptDirectoriesEmpty(), TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 144 - "verifyReleaseWithKey"
Cohesion: 0.40
Nodes (5): parseChecksum(), parsePublicKey(), TestParseChecksum(), verifyRelease(), verifyReleaseWithKey()

### Community 145 - "io.Reader"
Cohesion: 0.33
Nodes (4): BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype(), progressReader

### Community 146 - "Guard"
Cohesion: 0.07
Nodes (13): deferredBinary, rootHandle, openRoot(), SharedPinDegraded(), canonicalBinaryPath(), confinedEntry(), eventMask(), Guard (+5 more)

### Community 147 - "tempRow"
Cohesion: 0.29
Nodes (4): tempGrant, tempMaskOp, tempRow, liveTrust()

### Community 149 - "configedit_test.go"
Cohesion: 0.57
Nodes (7): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession()

### Community 150 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 151 - "bpfstats/main.go"
Cohesion: 0.19
Nodes (16): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), inertProgram(), kernelHasFunc(), TestKernelHasFunc(), trustSpec() (+8 more)

### Community 152 - "lockRootRecovering"
Cohesion: 0.29
Nodes (8): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), PinPrefix(), sanitizeGen(), SharedPinPrefix(), TestPinPrefixDistinctAndStable(), TestPinPrefixShape()

### Community 153 - "newTestSession"
Cohesion: 0.50
Nodes (3): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity()

### Community 155 - "find.go"
Cohesion: 0.12
Nodes (24): changelogText(), newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes(), findButton (+16 more)

### Community 156 - "prepareDaemonStart"
Cohesion: 0.11
Nodes (21): CheckBPFLSM(), CheckEBPF(), ParseNetEventsFlag(), configureDaemonLogging(), loadDaemonConfig(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath() (+13 more)

### Community 157 - "VerifyInstalledResourcesLocked"
Cohesion: 0.67
Nodes (3): TestVerifyResourcesLocked(), VerifyInstalledResourcesLocked(), VerifyResourcesLocked()

### Community 165 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 166 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 172 - "main_test.go"
Cohesion: 0.17
Nodes (6): infraContainerPaths(), tailLast(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

## Knowledge Gaps
- **53 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 375 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **32 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `testcontainers.Container`, `main_test.go`, `IntegrationSuite`, `shQuote`, `go_pkg_github_com_testcontainers_testcontainers_go`, `.startGuardStd`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.059) - this node is a cross-community bridge._
- **Why does `Monitor` connect `Monitor` to `daemon/daemon.go`, `NetEventType`, `sync.Mutex`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `GuardInodeKey`, `DaemonEvent`, `binaryVerifyState`, `GuardEvent`, `runDaemon`, `fakeDaemon`, `sync.Mutex`, `time.Time`, `os.File`, `BinaryEntry`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _53 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1044482673875645 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.024 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.12317073170731707 - nodes in this community are weakly interconnected._