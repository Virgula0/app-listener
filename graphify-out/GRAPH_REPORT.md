# Graph Report - app-listener  (2026-10-08)

## Corpus Check
- 323 files · ~1,504,562 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4286 nodes · 17292 edges · 170 communities (143 shown, 27 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1258 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6e351d28`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- .Unlock
- daemonconfig_test.go
- Vault
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- auth.go
- go_pkg_strings
- buildTrustedSet
- runUpdate
- install.go
- os.File
- IntegrationSuite
- guard_trust.bpf.c
- confgen.go
- fileEditModel
- BinaryEntry
- usecase/daemon_test.go
- Resource
- KernelDev
- Candidate
- guard.bpf.c
- Run
- wipe
- runDaemon
- buildOneGuard
- New
- serve.go
- engine
- IntegrationSuite
- fileedit_model_test.go
- DaemonEvent
- daemonconfig.go
- codeedit.go
- GuardEvent
- bpfstats/main.go
- highlight_test.go
- NetworkMonitor
- absPath
- IntegrationSuite
- IntegrationSuite
- ListUsers
- Monitor
- .syncSetsLocked
- IntegrationSuite
- CandidateDir
- check_and_emit_args
- IntegrationSuite
- github.com/cilium/ebpf.Map
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- GuardInodeKey
- controlServer
- SelectOrphans
- EventType
- TrustGuard
- runNetworkGuard
- exe_supersede.h
- fakeGuardRepo
- string
- pinstate.go
- SanitizeText
- runMonitor
- github.com/spf13/cobra.Command
- catalogRefresher
- copyRegularAt
- netModel
- audit.go
- btrfsLayoutFrom
- Ledger
- parseGuardEvent
- planUpdaters
- process_vm_readv.c
- Read
- go_pkg_os
- sanitizeTerminalText
- ParseEventsFlag
- .step
- editconfig.go
- render.go
- time.Time
- .GrantTemporaryAccess
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- refresh.go
- NetworkMonitorUseCase
- netmc.c
- backups.go
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- RawTarget
- TestMonitorUseCaseLifecycle
- IntegrationSuite
- Well-known bypass classes
- newTestVetter
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- bottomBar
- Build & sign release assets (reusable workflow)
- Highlighter
- sync.Mutex
- go_pkg_testing
- launchKey
- app-listener Terminal Demo Recording
- Config
- Vault
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- bufio.Reader
- TestMain
- WebSocket /ws client
- NewDumpHook
- launch_probe.c
- install.sh
- go_pkg_github_com_testcontainers_testcontainers_go
- launch_rule_at
- make test-integration (rootful Docker suite)
- vaultFS
- supersede_probe.c
- inode_dev
- trace-app-libs.sh
- IntegrationSuite
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- netGuardModel
- time.Duration
- net_tester/main.go
- update.go
- Guard
- tempRow
- runForward
- configedit_test.go
- IntegrationSuite
- io.Reader
- bpfLsmListed
- update_test.go
- find.go
- lockAndDeprovision
- Guard
- find_test.go
- classifySupportError
- loadDaemonConfig
- graphify-refresh.sh
- mustSubkey
- .runEditorHarness
- newChangelogModel
- downloadAndVerify
- fakeDaemon
- ResolvePinBase
- Vault
- backingDeviceUnion

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

## Communities (170 total, 27 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (52): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ldConfParser, multicallCacheKey (+44 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (116): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+108 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - ".Unlock"
Cohesion: 0.19
Nodes (9): isFileVaultCiphertext(), isRegularFileTarget(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall(), TestNewBoundedKeyFnRetryAborts() (+1 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (79): TestSetSectionWhitelistPreservesGroupStructure(), updateCatalogConfig(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines() (+71 more)

### Community 5 - "Vault"
Cohesion: 0.14
Nodes (19): lockOneRoot(), encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption() (+11 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.06
Nodes (70): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), emit_event(), emit_event_kern() (+62 more)

### Community 7 - "shellrc.go"
Cohesion: 0.05
Nodes (82): trustManager, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), setupBunTmpdir(), deploy(), TestDaemonUnitWithMetadataOutput(), addKeysToAgent() (+74 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (23): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf() (+15 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "auth.go"
Cohesion: 0.17
Nodes (25): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+17 more)

### Community 11 - "go_pkg_strings"
Cohesion: 0.30
Nodes (5): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed(), hint

### Community 12 - "buildTrustedSet"
Cohesion: 0.14
Nodes (17): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+9 more)

### Community 13 - "runUpdate"
Cohesion: 0.13
Nodes (22): changelogText(), showChangelog(), TestChangelogText(), confirmUpdate(), fetchReleases(), filterChannel(), isTerminal(), latestTime() (+14 more)

### Community 14 - "install.go"
Cohesion: 0.07
Nodes (42): confirmMasterKeyOverwrite(), runGenKey(), reloadDaemonForPasswordChange(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey() (+34 more)

### Community 15 - "os.File"
Cohesion: 0.05
Nodes (58): hashNow(), BinaryStat, Multicall, MulticallKind, TrustGuard, ExeKey(), Guard, IsMulticallRefusal() (+50 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "confgen.go"
Cohesion: 0.09
Nodes (41): TestDiffMergePreservesExistingAndParses(), sectionsFromCandidates(), LibraryBlock, Section, ConfSafePath(), findSectionEnd(), findSectionStart(), firstContentLine() (+33 more)

### Community 19 - "fileEditModel"
Cohesion: 0.11
Nodes (10): clipLabel(), humanSize(), listEntries(), sortNodes(), chownUser, dirEntry, fileEditMode, fileEditModel (+2 more)

### Community 20 - "BinaryEntry"
Cohesion: 0.08
Nodes (19): deferredBinary, VettedInode, BinaryEntry, canonicalBinaryPath(), canonicalPaths(), commMatchesGuardedBinary(), ConfinedEntry(), confinedEntry() (+11 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.13
Nodes (54): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+46 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (10): ownerPaths(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "KernelDev"
Cohesion: 0.19
Nodes (19): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), openRoot(), physicalParent(), statFD() (+11 more)

### Community 24 - "Candidate"
Cohesion: 0.14
Nodes (21): appendSectionsAndEdit(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+13 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (52): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class() (+44 more)

### Community 26 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 27 - "wipe"
Cohesion: 0.23
Nodes (18): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+10 more)

### Community 28 - "runDaemon"
Cohesion: 0.09
Nodes (30): CheckBPFLSM(), buildConcurrency(), buildGuards(), configureDaemonLogging(), drainEphemeral(), makeReloadHandler(), newPinGeneration(), notifySystemdReady() (+22 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.13
Nodes (23): buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions(), newSelfGuards(), selfGuards, AdmissionCheck (+15 more)

### Community 30 - "New"
Cohesion: 0.22
Nodes (15): harnessSuite, looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip() (+7 more)

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "engine"
Cohesion: 0.09
Nodes (7): TestTieRank(), engine, pathWithin(), tieRank(), deleteInoKeys(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent()

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.08
Nodes (47): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel() (+39 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.11
Nodes (32): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+24 more)

### Community 36 - "daemonconfig.go"
Cohesion: 0.15
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 37 - "codeedit.go"
Cohesion: 0.25
Nodes (17): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+9 more)

### Community 38 - "GuardEvent"
Cohesion: 0.08
Nodes (19): selfProtectSpecs(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), selfProtectSpec (+11 more)

### Community 39 - "bpfstats/main.go"
Cohesion: 0.19
Nodes (16): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), inertProgram(), kernelHasFunc(), TestKernelHasFunc(), trustSpec() (+8 more)

### Community 40 - "highlight_test.go"
Cohesion: 0.18
Nodes (23): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+15 more)

### Community 41 - "NetworkMonitor"
Cohesion: 0.10
Nodes (13): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), NewNetworkMonitor() (+5 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.23
Nodes (7): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "ListUsers"
Cohesion: 0.13
Nodes (23): askBunTmpdirUsers(), gatherPerUserSetup(), applyDiffAdditions(), collectDiffAdditions(), configPaths(), pathCovered(), restoreDaemonAfterDiffAbort(), runDiffCatalog() (+15 more)

### Community 46 - "Monitor"
Cohesion: 0.06
Nodes (8): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), Monitor, monitorUnitTest, pathCache, pathCacheEntry

### Community 48 - ".syncSetsLocked"
Cohesion: 0.11
Nodes (18): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, intersectAllows(), isSubset(), planMemberRows() (+10 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.17
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "CandidateDir"
Cohesion: 0.06
Nodes (33): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, GlobChild, statFunc, compileGlobNames() (+25 more)

### Community 51 - "check_and_emit_args"
Cohesion: 0.18
Nodes (14): exe_is_superseded(), exe_refused(), check_and_emit_args(), code_suspect(), current_is_inspector(), exe_tag_bits(), fill_path(), get_current_exe_inode() (+6 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (16): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+8 more)

### Community 53 - "github.com/cilium/ebpf.Map"
Cohesion: 0.09
Nodes (16): BtrfsLayout, checkBtrfsKeys(), ensureBtrfsLayout(), syncU32Map(), loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows(), syncMap() (+8 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.11
Nodes (7): clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 56 - "GuardInodeKey"
Cohesion: 0.07
Nodes (25): refusedReplacements, engine, Inspector, SetInspectors(), engine, supersedeUnnamed(), appletTag(), buildOf() (+17 more)

### Community 57 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 58 - "SelectOrphans"
Cohesion: 0.23
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 59 - "EventType"
Cohesion: 0.11
Nodes (10): BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryEventsReadOnlyRejectsRestriction(), Run(), Cstr(), DecodeBpfEvent() (+2 more)

### Community 60 - "TrustGuard"
Cohesion: 0.12
Nodes (8): trustHook, trustRows, setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied(), NewTrustGuard(), TestTrustGuardStartFailsWhenLibraryAllowlistHookCannotAttach()

### Community 61 - "runNetworkGuard"
Cohesion: 0.16
Nodes (12): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode, NetworkGuardRepository (+4 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.20
Nodes (16): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+8 more)

### Community 64 - "string"
Cohesion: 0.14
Nodes (9): denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 65 - "pinstate.go"
Cohesion: 0.12
Nodes (31): lockRootRecovering(), relockStaleVaults(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState() (+23 more)

### Community 66 - "SanitizeText"
Cohesion: 0.24
Nodes (6): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), binaryVetter, vetDecision, SanitizeText()

### Community 67 - "runMonitor"
Cohesion: 0.11
Nodes (17): CheckEBPF(), MakeDisplayPaths(), ParseNetEventsFlag(), runHeadless(), runMonitor(), runTUI(), runNetworkMonitor(), runTUI() (+9 more)

### Community 68 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 69 - "catalogRefresher"
Cohesion: 0.17
Nodes (8): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), catalogRefresher, Parse(), validateResources()

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 73 - "audit.go"
Cohesion: 0.23
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.19
Nodes (8): fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers(), TestBtrfsLayoutRejectsUnexpectedShapes()

### Community 75 - "Ledger"
Cohesion: 0.09
Nodes (19): grantNewBinaries(), candidates(), confirmable(), confirmAll(), describe(), run(), EnsurePlaceholders(), Ledger (+11 more)

### Community 76 - "parseGuardEvent"
Cohesion: 0.22
Nodes (7): bpfGuardEvent, fsGateLabel(), parseGuardEvent(), processGateLabel(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels(), suspectLabel()

### Community 78 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "Read"
Cohesion: 0.10
Nodes (20): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+12 more)

### Community 81 - "go_pkg_os"
Cohesion: 0.06
Nodes (41): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), decryptStep(), pickDirsToDecrypt(), revertSystemFiles(), checkKeyLen() (+33 more)

### Community 82 - "sanitizeTerminalText"
Cohesion: 0.17
Nodes (8): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), TestSanitizeTerminalText()

### Community 83 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 84 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 85 - "editconfig.go"
Cohesion: 0.26
Nodes (12): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+4 more)

### Community 86 - "render.go"
Cohesion: 0.15
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "time.Time"
Cohesion: 0.14
Nodes (14): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+6 more)

### Community 88 - ".GrantTemporaryAccess"
Cohesion: 0.14
Nodes (12): TemporaryGrant, TempRule, applyTempGrants(), daemonUseCase, TemporaryAccess, logTempGrant(), newExeTap(), planTempGrants() (+4 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "refresh.go"
Cohesion: 0.21
Nodes (14): logRefreshChange(), TestLogRefreshChange(), RefreshOptions, SectionScan, IsLibraryDirective(), admitNew(), SectionChange, keepExisting() (+6 more)

### Community 92 - "NetworkMonitorUseCase"
Cohesion: 0.15
Nodes (5): runHeadless(), runHeadless(), ProtocolString(), NetworkMonitorRepository, NetworkMonitorUseCase

### Community 93 - "netmc.c"
Cohesion: 0.12
Nodes (3): do_connect(), errname(), main()

### Community 94 - "backups.go"
Cohesion: 0.13
Nodes (24): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir(), Delete() (+16 more)

### Community 95 - "NetGuard"
Cohesion: 0.13
Nodes (8): foreignPinnedLink(), eventsetSummary(), eventTypeKey(), mcNameKey(), modeLabel(), NewNetGuard(), statInodeKey(), NetGuard

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
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 100 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (11): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

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

### Community 107 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "Highlighter"
Cohesion: 0.12
Nodes (7): span, diffModel, editorModel, Highlighter, lexerFor(), shebangLexer(), tokenClassOf()

### Community 110 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 111 - "go_pkg_testing"
Cohesion: 0.09
Nodes (13): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+5 more)

### Community 112 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "Config"
Cohesion: 0.28
Nodes (10): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager, startTrustGuard(), trustStartupError(), bunEntryConfigured(), bunLaunchersFor(), inspectorAdmitter (+2 more)

### Community 115 - "Vault"
Cohesion: 0.23
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "bufio.Reader"
Cohesion: 0.16
Nodes (13): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+5 more)

### Community 121 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

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

### Community 127 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "vaultFS"
Cohesion: 0.16
Nodes (13): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+5 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "inode_dev"
Cohesion: 0.27
Nodes (7): ctx_ptr(), inode_dev(), sb_dev(), guard_exec_applet(), guard_inode_free(), mc_basename(), trust_inode_free()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "netGuardModel"
Cohesion: 0.24
Nodes (5): NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 143 - "time.Duration"
Cohesion: 0.15
Nodes (6): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), parseErr(), debouncer, liveSession

### Community 144 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 145 - "update.go"
Cohesion: 0.22
Nodes (7): compareStableVersions(), newerThanStable(), parseStableVersion(), parseVerPart(), TestNewerThanStable(), TestParseStableVersion(), stableVersion

### Community 146 - "Guard"
Cohesion: 0.08
Nodes (13): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, rootHandle, sharedHashEntry, SharedPinDegraded() (+5 more)

### Community 147 - "tempRow"
Cohesion: 0.29
Nodes (4): tempGrant, tempMaskOp, tempRow, liveTrust()

### Community 148 - "runForward"
Cohesion: 0.18
Nodes (13): chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents(), pingInterval() (+5 more)

### Community 149 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 151 - "io.Reader"
Cohesion: 0.22
Nodes (7): entryOf(), BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype(), progressReader

### Community 152 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 153 - "update_test.go"
Cohesion: 0.15
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 154 - "find.go"
Cohesion: 0.20
Nodes (7): findButton, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 155 - "lockAndDeprovision"
Cohesion: 0.21
Nodes (8): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision(), modifiedContextWithSource()

### Community 156 - "Guard"
Cohesion: 0.30
Nodes (3): exeRow, Guard, planTempBlock()

### Community 157 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 158 - "classifySupportError"
Cohesion: 0.18
Nodes (10): classifySetupError(), classifySupportError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), TestClassifySupportErrorEncryptionNotEnabled(), TestClassifySupportErrorEncryptionNotEnabledF2fs(), TestClassifySupportErrorGeneric() (+2 more)

### Community 159 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), resolveConfigPath(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), SetUserHomes(), TestOpenSystemPlacedRefusesUserHome()

### Community 161 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 163 - "newChangelogModel"
Cohesion: 0.31
Nodes (6): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 164 - "downloadAndVerify"
Cohesion: 0.31
Nodes (7): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), sanityCheckBinary(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), updateHTTPClient()

### Community 166 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

## Knowledge Gaps
- **53 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 381 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **27 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `go_pkg_os`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`?**
  _High betweenness centrality (0.069) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `fakeDaemon`, `GuardEvent`, `sync.Mutex`, `os.File`, `BinaryEntry`, `time.Time`, `GuardInodeKey`, `runDaemon`, `buildOneGuard`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **Why does `vaultFS` connect `vaultFS` to `go_pkg_os`, `fileedit_model_test.go`, `fileEditModel`, `os.File`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _53 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12186507936507937 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02577980344680245 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._