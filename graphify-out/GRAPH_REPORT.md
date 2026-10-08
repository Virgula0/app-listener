# Graph Report - app-listener  (2026-10-08)

## Corpus Check
- 323 files · ~1,504,995 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4289 nodes · 17300 edges · 174 communities (144 shown, 30 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1262 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1f42e082`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- secureResources
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- auth.go
- uninstall_system.go
- buildTrustedSet
- update.go
- install.go
- ExeKey
- IntegrationSuite
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- BinaryEntry
- usecase/daemon_test.go
- Resource
- inodeChain
- Candidate
- guard.bpf.c
- dataRowRenderer
- Monitor
- runDaemon
- buildOneGuard
- filevault_test.go
- serve.go
- engine
- IntegrationSuite
- fileedit_model_test.go
- DaemonEvent
- watchGroup
- serveWebSocket
- GuardEvent
- bpfstats/main.go
- github.com/charmbracelet/bubbles/textarea.Model
- NetworkMonitor
- absPath
- filevault.go
- IntegrationSuite
- update_test.go
- monitorUnitTest
- github.com/cilium/ebpf.Map
- IntegrationSuite
- CandidateDir
- check_and_emit_args
- IntegrationSuite
- syncMap
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- GuardInodeKey
- controlServer
- fscrypt.go
- EventType
- TrustGuard
- LibraryClosure
- exe_supersede.h
- fakeGuardRepo
- string
- pinstate.go
- time.Time
- github.com/spf13/cobra.Command
- headerWidget
- catalogRefresher
- copyRegularAt
- netModel
- Config
- btrfsLayoutFrom
- Ledger
- SanitizeText
- planUpdaters
- process_vm_readv.c
- Read
- offline.go
- .SetGlobReservations
- NewEventFanout
- OpenSystemPlaced
- updateCatalogConfig
- render.go
- watchSet
- TempBinary
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- supersedeMaps
- Vault
- netmc.c
- liveSession
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- RawTarget
- eventUnitTest
- IntegrationSuite
- Well-known bypass classes
- newTestVetter
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- .AllowReplacement
- Build & sign release assets (reusable workflow)
- ClassifyMulticall
- sync.Mutex
- writeWithin
- runGuard
- app-listener Terminal Demo Recording
- NetEvent
- ResolveLibraryClosure
- fileedit_symlink_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- bufio.Reader
- watchPattern
- WebSocket /ws client
- NewDumpHook
- launch_probe.c
- install.sh
- IntegrationSuite
- launch_rule_at
- make test-integration (rootful Docker suite)
- os.File
- supersede_probe.c
- inode_dev
- trace-app-libs.sh
- IntegrationSuite
- openBunTmp
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- netGuardModel
- time.Duration
- net_tester/main.go
- testLimiter
- Guard
- tempRow
- bottomBar
- configedit_test.go
- downloadFile
- findTopLevelWindow
- bpfLsmListed
- StripPinnedTempAllows
- .readEvent
- launchKey
- IntegrationSuite
- .pathInfo
- validateFlags
- loadDaemonConfig
- graphify-refresh.sh
- mustSubkey
- ResolvePinBase
- newChangelogModel
- Vault
- fakeDaemon
- TestMain
- .runEditorHarness
- .exeKeys
- parseResize
- formatEntry
- ldConfParser
- CleanupStalePins

## God Nodes (most connected - your core abstractions)
1. `Config` - 93 edges
2. `Guard` - 93 edges
3. `GuardInodeKey` - 86 edges
4. `Load()` - 85 edges
5. `absPath()` - 77 edges
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

## Communities (174 total, 30 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (49): Execute(), tempJournalRow, classCacheKey, fdStat, launchRule, MulticallError, trustEvent, BinaryStat (+41 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (123): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+115 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.13
Nodes (18): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+10 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (74): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+66 more)

### Community 5 - "secureResources"
Cohesion: 0.13
Nodes (14): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse() (+6 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.06
Nodes (70): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), emit_event(), emit_event_kern() (+62 more)

### Community 7 - "shellrc.go"
Cohesion: 0.05
Nodes (95): trustManager, askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), restoreDaemonAfterDiffAbort() (+87 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+13 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "auth.go"
Cohesion: 0.08
Nodes (36): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+28 more)

### Community 12 - "buildTrustedSet"
Cohesion: 0.16
Nodes (16): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+8 more)

### Community 13 - "update.go"
Cohesion: 0.12
Nodes (29): changelogText(), showChangelog(), TestChangelogText(), applyUpdate(), compareStableVersions(), confirmUpdate(), downloadAndVerify(), filterChannel() (+21 more)

### Community 14 - "install.go"
Cohesion: 0.08
Nodes (47): runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), installBinaryOnly(), parseMaintenanceFlags(), preflightDeployedBinary() (+39 more)

### Community 15 - "ExeKey"
Cohesion: 0.12
Nodes (26): confirmUserPlaced(), TrustGuard, ExeKey(), IsMulticallRefusal(), MulticallRefusal(), RefuseMulticallLines(), registerBuild(), siblingNames() (+18 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "GenerateConf"
Cohesion: 0.06
Nodes (51): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, Section, findSectionEnd(), findSectionStart(), firstContentLine() (+43 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "BinaryEntry"
Cohesion: 0.09
Nodes (17): deferredBinary, VettedInode, BinaryEntry, BinariesSummary(), canonicalBinaryPath(), ConfinedEntry(), confinedEntry(), eventMask() (+9 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.11
Nodes (54): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+46 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (11): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), uniqueEncryptionRoots() (+3 more)

### Community 23 - "inodeChain"
Cohesion: 0.44
Nodes (9): inodeChain(), isProcFD(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow(), TestInodeChainIsPhysicalThroughSymlink(), TestInodeChainOfPinnedFileRoot() (+1 more)

### Community 24 - "Candidate"
Cohesion: 0.11
Nodes (27): appendSectionsAndEdit(), configPaths(), pathCovered(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories() (+19 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (52): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class() (+44 more)

### Community 28 - "runDaemon"
Cohesion: 0.09
Nodes (31): CheckBPFLSM(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), makeReloadHandler(), newPinGeneration(), notifySystemdReady() (+23 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.10
Nodes (31): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan() (+23 more)

### Community 30 - "filevault_test.go"
Cohesion: 0.24
Nodes (19): decryptDirectories(), TestDecryptDirectoriesEmpty(), inode(), TestEnsureRecoverySidecarPlaceholder(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceRecoversFromInterruptedTransform() (+11 more)

### Community 31 - "serve.go"
Cohesion: 0.14
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

### Community 32 - "engine"
Cohesion: 0.09
Nodes (6): TestTieRank(), engine, pathWithin(), tieRank(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent()

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.10
Nodes (37): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel() (+29 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.11
Nodes (25): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+17 more)

### Community 36 - "watchGroup"
Cohesion: 0.09
Nodes (38): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+30 more)

### Community 37 - "serveWebSocket"
Cohesion: 0.22
Nodes (8): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 38 - "GuardEvent"
Cohesion: 0.07
Nodes (18): bpfGuardEvent, commMatchesGuardedBinary(), fsGateLabel(), GuardEvent, logBacklogDenial(), parseGuardEvent(), processGateLabel(), TestParseGuardEventFsGateLabels() (+10 more)

### Community 39 - "bpfstats/main.go"
Cohesion: 0.19
Nodes (16): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), inertProgram(), kernelHasFunc(), TestKernelHasFunc(), trustSpec() (+8 more)

### Community 40 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.05
Nodes (64): class, findButton, span, diffModel, editorModel, classOf(), column(), dedent() (+56 more)

### Community 41 - "NetworkMonitor"
Cohesion: 0.18
Nodes (5): NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "filevault.go"
Cohesion: 0.24
Nodes (15): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, looksLikeFileVaultRecord(), openFileVaultWithMasterKey() (+7 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.23
Nodes (7): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "update_test.go"
Cohesion: 0.11
Nodes (13): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), testKeyPair() (+5 more)

### Community 46 - "monitorUnitTest"
Cohesion: 0.10
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 48 - "github.com/cilium/ebpf.Map"
Cohesion: 0.11
Nodes (19): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), deleteResKeys(), engine, isSubset(), planMemberRows() (+11 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.17
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "CandidateDir"
Cohesion: 0.07
Nodes (25): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), confSafeMatch() (+17 more)

### Community 51 - "check_and_emit_args"
Cohesion: 0.18
Nodes (14): exe_is_superseded(), exe_refused(), check_and_emit_args(), code_suspect(), current_is_inspector(), exe_tag_bits(), fill_path(), get_current_exe_inode() (+6 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.10
Nodes (8): clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel, changelogModel

### Community 56 - "GuardInodeKey"
Cohesion: 0.07
Nodes (24): refusedReplacements, TestNearestRootClaims(), nearestRootClaims(), deleteInoKeys(), engine, Inspector, SetInspectors(), engine (+16 more)

### Community 57 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 58 - "fscrypt.go"
Cohesion: 0.08
Nodes (38): confirmMasterKeyOverwrite(), runGenKey(), ensureMasterKey(), cleanOrphanedMetadata(), checkKeyLen(), classifySetupError(), GenerateMasterKey(), isNotEncryptedErr() (+30 more)

### Community 59 - "EventType"
Cohesion: 0.09
Nodes (14): systemPatterns(), systemRule(), BpfEvent, EventType, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+6 more)

### Community 60 - "TrustGuard"
Cohesion: 0.13
Nodes (6): trustHook, trustRows, setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 61 - "LibraryClosure"
Cohesion: 0.15
Nodes (13): defaultLibDirs(), elfInterp(), isScript(), LibraryClosure(), rootOwnedInode(), rootOwnedSafe(), TestResolveLibraryClosureRejectsUserWritableOriginDir(), runsNatively() (+5 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.20
Nodes (16): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+8 more)

### Community 64 - "string"
Cohesion: 0.14
Nodes (9): denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 65 - "pinstate.go"
Cohesion: 0.13
Nodes (28): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive() (+20 more)

### Community 66 - "time.Time"
Cohesion: 0.25
Nodes (3): guiModel, newColumnState(), Run()

### Community 67 - "github.com/spf13/cobra.Command"
Cohesion: 0.06
Nodes (43): AddServeFlags(), CheckEBPF(), credentialFlagState(), isLoopbackHost(), ParseEventsFlag(), ParseNetEventsFlag(), ParseServeFlags(), requestedBoolFlag() (+35 more)

### Community 68 - "headerWidget"
Cohesion: 0.20
Nodes (5): columnState, dataRow, headerWidget, newDataRow(), newHeaderWidget()

### Community 69 - "catalogRefresher"
Cohesion: 0.10
Nodes (16): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), awaitStartupOrSignal(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+8 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 73 - "Config"
Cohesion: 0.11
Nodes (29): ServeConfig, newDaemonModel(), runDaemonUI(), runServedTUI(), runTUI(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager (+21 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "Ledger"
Cohesion: 0.10
Nodes (16): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), EnsurePlaceholders(), Ledger (+8 more)

### Community 76 - "SanitizeText"
Cohesion: 0.17
Nodes (10): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), runHeadless(), runHeadless(), binaryVetter, vetDecision (+2 more)

### Community 78 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "Read"
Cohesion: 0.10
Nodes (19): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+11 more)

### Community 81 - "offline.go"
Cohesion: 0.06
Nodes (44): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), deletePostBackups(), restoreBackups(), askDecrypt(), decryptStep(), offerBackupCleanup() (+36 more)

### Community 82 - ".SetGlobReservations"
Cohesion: 0.20
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), TrustGuard, resolveBits(), resolveChildren(), statKey() (+1 more)

### Community 83 - "NewEventFanout"
Cohesion: 0.22
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 84 - "OpenSystemPlaced"
Cohesion: 0.13
Nodes (18): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), OpenSystemPlaced() (+10 more)

### Community 85 - "updateCatalogConfig"
Cohesion: 0.16
Nodes (16): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+8 more)

### Community 86 - "render.go"
Cohesion: 0.14
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "watchSet"
Cohesion: 0.24
Nodes (6): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet

### Community 88 - "TempBinary"
Cohesion: 0.14
Nodes (15): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess, logTempGrant() (+7 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "supersedeMaps"
Cohesion: 0.18
Nodes (12): checkBtrfsKeys(), ensureBtrfsLayout(), dropSupersededLocked(), PruneSuperseded(), supersedeMaps(), supersedePruneLoop(), warnLostStamps(), NewTrustGuard() (+4 more)

### Community 92 - "Vault"
Cohesion: 0.23
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 93 - "netmc.c"
Cohesion: 0.12
Nodes (3): do_connect(), errname(), main()

### Community 95 - "NetGuard"
Cohesion: 0.14
Nodes (9): NetEventType, foreignPinnedLink(), NetEventTypes(), ParseNetEventType(), eventsetSummary(), eventTypeKey(), modeLabel(), NewNetGuard() (+1 more)

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

### Community 100 - "eventUnitTest"
Cohesion: 0.18
Nodes (3): eventUnitTest, DecodeBpfEvent(), editorHarness

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

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "ClassifyMulticall"
Cohesion: 0.22
Nodes (11): Multicall, MulticallKind, multicallFamily(), singleBinaryPaths(), ClassifyMulticall(), familyOf(), listUutilsApplets(), scanMulticallFamily() (+3 more)

### Community 110 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 111 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 112 - "runGuard"
Cohesion: 0.14
Nodes (15): selfProtectSpecs(), displayEntries(), holdForward(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless() (+7 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "NetEvent"
Cohesion: 0.13
Nodes (6): NetEvent, NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase(), newFakeNetworkMonitorRepo(), fakeNetworkMonitorRepo

### Community 115 - "ResolveLibraryClosure"
Cohesion: 0.18
Nodes (11): hostView(), mountVouchesOwnership(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch() (+3 more)

### Community 116 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 121 - "watchPattern"
Cohesion: 0.39
Nodes (6): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern

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

### Community 127 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "os.File"
Cohesion: 0.13
Nodes (15): checkFreshDir(), openDirNoFollow(), renewDir(), Guard, genericElectron(), runtimeClass(), scanRuntimeClass(), TestRuntimeClassGenericElectron() (+7 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "inode_dev"
Cohesion: 0.27
Nodes (7): ctx_ptr(), inode_dev(), sb_dev(), guard_exec_applet(), guard_inode_free(), mc_basename(), trust_inode_free()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "netGuardModel"
Cohesion: 0.10
Nodes (17): NetGuardEvent, NetworkGuardRepository, protoString(), NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle() (+9 more)

### Community 143 - "time.Duration"
Cohesion: 0.19
Nodes (5): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), pingInterval(), debouncer

### Community 144 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 145 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 146 - "Guard"
Cohesion: 0.08
Nodes (7): binaryVerifyState, rootHandle, sharedHashEntry, openRoot(), SharedPinDegraded(), Guard, hashBinaryShared()

### Community 147 - "tempRow"
Cohesion: 0.17
Nodes (7): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust(), planTempBlock()

### Community 148 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 149 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 150 - "downloadFile"
Cohesion: 0.28
Nodes (6): downloadFile(), downloadReleaseFiles(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), updateHTTPClient(), progressReader

### Community 151 - "findTopLevelWindow"
Cohesion: 0.38
Nodes (4): findTopLevelWindow(), floatWindow(), internAtom(), setFloatHint()

### Community 152 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 153 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 154 - ".readEvent"
Cohesion: 0.33
Nodes (4): NetBpfEvent, Cstr(), FormatAddr(), Ntohs()

### Community 155 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 157 - ".pathInfo"
Cohesion: 0.48
Nodes (3): evalSymlinksOrEmpty(), pathCache, pathCacheEntry

### Community 158 - "validateFlags"
Cohesion: 0.50
Nodes (4): TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags()

### Community 159 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), resolveConfigPath(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), SetUserHomes(), TestOpenSystemPlacedRefusesUserHome()

### Community 161 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 162 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 163 - "newChangelogModel"
Cohesion: 0.53
Nodes (5): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes()

### Community 166 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 169 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 170 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 172 - "ldConfParser"
Cohesion: 0.50
Nodes (3): ldConfParser, ldConfInclude(), ldSoConfDirs()

### Community 173 - "CleanupStalePins"
Cohesion: 0.40
Nodes (5): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

## Knowledge Gaps
- **53 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 381 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **30 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `auth.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.065) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `os.File`, `time.Time`, `DaemonEvent`, `fakeDaemon`, `GuardEvent`, `sync.Mutex`, `runGuard`, `BinaryEntry`, `GuardInodeKey`, `runDaemon`, `buildOneGuard`?**
  _High betweenness centrality (0.058) - this node is a cross-community bridge._
- **Why does `classifyMigrationTarget()` connect `Vault` to `uninstall_system.go`, `filevault.go`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _53 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09669547534316218 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.024 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._